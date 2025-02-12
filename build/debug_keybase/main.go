package main

import (
	"bytes"
	"crypto/md5" // nolint:gosec // Weak hashing function only used to check if the file has been changed
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/korovkin/limiter"
	"gopkg.in/yaml.v2"

	"github.com/pokt-network/pocket/app/client/keybase"
	"github.com/pokt-network/pocket/runtime/configs"
	"github.com/pokt-network/pocket/runtime/defaults"
	cryptoPocket "github.com/pokt-network/pocket/shared/crypto"
	"github.com/pokt-network/pocket/shared/utils"
)

const (
	// NOTE: This is the number of validators in the private-keys.yaml manifest file
	numValidators = 999

	// Increasing this number is linearly proportional to the amount of RAM required for the debug client to start and import
	// pre-generated keys into the keybase. Beware that might cause OOM and process can exit with 137 status code.
	// 4 threads takes 350-400MiB from a few tests which sounds acceptable.
	debugKeybaseImportConcurrencyLimit = 4
)

type K8sSecret struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	MetaData   map[string]string `yaml:"metadata"`
	Type       string            `yaml:"type"`
	StringData map[string]string `yaml:"stringData"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <source_yaml> <target_folder>")
		return
	}
	sourceYamlPath := os.Args[1]
	targetFolderPath := os.Args[2]

	privateKeysYamlBytes, err := os.ReadFile(sourceYamlPath)
	if err != nil {
		fmt.Printf("Error reading source_yaml: %v\n", err)
		return
	}
	sourceYamlHash := md5.Sum(privateKeysYamlBytes) // nolint:gosec // Weak hashing function only used to check if the file has been changed
	hashString := fmt.Sprintf("%x.md5", sourceYamlHash)
	hashFilePath := filepath.Join(targetFolderPath, hashString)
	targetFilePath := filepath.Join(targetFolderPath, "debug_keybase.bak")
	if exists, _ := utils.FileExists(hashFilePath); !exists {
		cleanupStaleFiles(targetFolderPath)
		dumpKeybase(privateKeysYamlBytes, targetFilePath)
		createHashFile(hashFilePath)
	} else {
		fmt.Println("✅ Keybase dump already exists and in sync with YAML file")
	}
}

func dumpKeybase(privateKeysYamlBytes []byte, targetFilePath string) {
	fmt.Println("⚙️  Initializing debug Keybase...")

	validatorKeysPairMap, err := parsePrivateKeysFromEmbeddedYaml(privateKeysYamlBytes)
	if err != nil {
		panic(err)
	}

	kb, err := keybase.NewKeybase(&configs.KeybaseConfig{
		FilePath: defaults.DefaultRootDirectory + "/keys",
	})
	if err != nil {
		panic(err)
	}
	db, err := kb.GetBadgerDB()
	if err != nil {
		panic(err)
	}

	// Add validator addresses if not present
	fmt.Println("✍️  Debug keybase initializing... Adding all the validator keys")

	// Use writebatch to speed up bulk insert
	wb := db.NewWriteBatch()

	// Create a channel to receive errors from goroutines
	errCh := make(chan error, numValidators)

	limit := limiter.NewConcurrencyLimiter(debugKeybaseImportConcurrencyLimit)
	for _, privHexString := range validatorKeysPairMap {
		if _, err := limit.Execute(func() {
			// Import the keys into the keybase with no passphrase or hint as these are for debug purposes
			keyPair, err := cryptoPocket.CreateNewKeyFromString(privHexString, "", "")
			if err != nil {
				errCh <- err
				return
			}

			// Use key address as key in DB
			addrKey := keyPair.GetAddressBytes()

			// Encode KeyPair into []byte for value
			keypairBz, err := keyPair.Marshal()
			if err != nil {
				errCh <- err
				return
			}
			if err := wb.Set(addrKey, keypairBz); err != nil {
				errCh <- err
				return
			}
		}); err != nil {
			panic(err)
		}
	}

	if err := limit.WaitAndClose(); err != nil {
		panic(err)
	}

	// Check if any goroutines returned an error
	select {
	case err := <-errCh:
		panic(err)
	default:
	}

	if err := wb.Flush(); err != nil {
		panic(err)
	}

	fmt.Println("✅ Keybase initialized!")

	fmt.Println("⚙️  Creating a dump of the Keybase...")
	backupFile, err := os.Create(targetFilePath)
	if err != nil {
		panic(err)
	}
	defer backupFile.Close()
	if _, err := db.Backup(backupFile, 0); err != nil {
		panic(err)
	}

	// Close DB connection
	if err := kb.Stop(); err != nil {
		panic(err)
	}

	fmt.Printf("✅ Keybase dumped in %s\n", targetFilePath)
}

func parsePrivateKeysFromEmbeddedYaml(privateKeysYamlBytes []byte) ([]string, error) {
	// Parse the YAML file and load into the config struct
	decoder := yaml.NewDecoder(bytes.NewReader(privateKeysYamlBytes))
	keysList := make([]string, 0)

	for {
		var secret K8sSecret

		if err := decoder.Decode(&secret); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		for _, privHexString := range secret.StringData {
			keysList = append(keysList, privHexString)
		}

	}

	return keysList, nil
}

func cleanupStaleFiles(targetFolderPath string) {
	fmt.Printf("🧹 Cleaning up stale backup files in in %s\n", targetFolderPath)

	if err := filepath.Walk(targetFolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}

		if !info.IsDir() && (filepath.Ext(path) == ".bak" || filepath.Ext(path) == ".md5") {
			if err := os.Remove(path); err != nil {
				panic(err)
			}
			fmt.Println("🚮 Deleted file:", path)
		}

		return nil
	}); err != nil {
		panic(err)
	}
}

func createHashFile(hashFilePath string) {
	fmt.Printf("🔖 Creating the MD5 hash file used to track consistency: %s\n", hashFilePath)

	file, err := os.Create(hashFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	if _, err := file.WriteString("This file is used to check if the keybase dump is in sync with the YAML file. Its name is the MD5 hash of the private_keys.yaml"); err != nil {
		panic(err)
	}
}
