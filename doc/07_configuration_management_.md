# Chapter 7: Configuration Management

In the previous chapter, [Asynchronous Task Processing (Worker)](06_asynchronous_task_processing__worker__.md), we saw how SimpleBank handles background tasks like sending emails using a separate worker process connected via Redis. Both our main API server and the worker process need to know things like the database address, the Redis address, and secret keys for tokens or emails.

Where should we store these settings? We _could_ write them directly into our Go code. For example:

```go
// BAD PRACTICE - DON'T DO THIS!
const databasePassword = "supersecretpassword123"
const serverAddress = "127.0.0.1:8080"
```

But this is a bad idea!

- **Insecurity:** Anyone who can see the code can see sensitive passwords or keys.
- **Inflexibility:** What if we want to run the application on a different computer (like a production server) that uses a different database password or listens on a different address? We'd have to change the code and recompile it every time!

We need a way to manage these external settings separately from our application code. This is called **Configuration Management**.

## The Application's Instruction Manual and Toolkit

Think of Configuration Management like your application's instruction manual and toolkit combined.

- **Instruction Manual:** It tells the application _where_ things are located (e.g., the address of the database, the address of the Redis server).
- **Toolkit:** It stores necessary tools and keys (e.g., the secret key needed to create secure [Authentication & Authorization (Token-Based)](05_authentication___authorization__token_based__.md) tokens, the password for sending emails).

Instead of hardcoding these details, the application reads them from a configuration source when it starts up. This makes it easy to change settings for different environments (like your local computer for development vs. a powerful server for production) without touching the actual Go code.

## Key Concepts: How SimpleBank Manages Settings

SimpleBank uses a popular Go library called **Viper** to handle configuration. Here's how it works:

1.  **Configuration Files (`app.env`):** We can store settings in a simple text file. By convention, these files often end in `.env` (for "environment"). Our file is named `app.env`. It looks like a list of `KEY=VALUE` pairs:

    ```env
    # Sample app.env file (simplified)

    ENVIRONMENT=development
    DB_SOURCE=postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable
    HTTP_SERVER_ADDRESS=0.0.0.0:8080
    GRPC_SERVER_ADDRESS=0.0.0.0:9090
    REDIS_ADDRESS=localhost:6379
    TOKEN_SYMMETRIC_KEY=12345678901234567890123456789012 # Must be 32 chars
    ACCESS_TOKEN_DURATION=15m
    REFRESH_TOKEN_DURATION=24h
    ```

    This file contains the database connection string, server addresses, the secret key for tokens, and how long tokens should last.

2.  **Environment Variables:** Settings can also be provided by the operating system's environment variables. This is very common in production deployments (e.g., using Docker or cloud platforms). Viper can automatically read these too! If an environment variable named `DB_SOURCE` exists, Viper can use its value, potentially overriding the value from the `app.env` file.

3.  **Viper Library:** Viper is the smart librarian that knows how to:

    - Find and read the `app.env` file.
    - Check for environment variables that match the keys in the file (e.g., `DB_SOURCE`).
    - Load these values into a structured format that our Go code can easily use.

4.  **`Config` Struct (`util/config.go`):** To make using these settings easy and type-safe in Go, we define a struct that mirrors the configuration keys.

    ```go
    // --- File: util/config.go (Simplified) ---
    package util

    import "time"

    // Config stores all configuration of the application.
    type Config struct {
    	Environment          string        `mapstructure:"ENVIRONMENT"`
    	DBSource             string        `mapstructure:"DB_SOURCE"`
    	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
    	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
    	RedisAddress         string        `mapstructure:"REDIS_ADDRESS"`
    	TokenSymmetricKey    string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
    	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
    	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
    	// ... other fields like email credentials, migration URL etc.
    }
    ```

    - Each field in the struct (`DBSource`, `HTTPServerAddress`, etc.) corresponds to a key in our `app.env` file or an environment variable.
    - The `mapstructure` tags tell Viper exactly which key maps to which field.

## Loading the Configuration

When the SimpleBank application starts (in `main.go`), one of the first things it does is load the configuration using a helper function in `util/config.go`.

```go
// --- File: util/config.go ---
package util

import (
	"time"
	"github.com/spf13/viper" // Import the Viper library
)

// ... (Config struct definition) ...

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (config Config, err error) {
	// Tell Viper where to look for the config file (the path argument)
	viper.AddConfigPath(path)
	// Tell Viper the name of the config file (without extension)
	viper.SetConfigName("app")
	// Tell Viper the type of the config file
	viper.SetConfigType("env") // e.g., "env", "json", "yaml"

	// Tell Viper to also read matching environment variables
	viper.AutomaticEnv()

	// Attempt to read the configuration file
	err = viper.ReadInConfig()
	if err != nil {
		// Handle error if the file wasn't found or couldn't be read
		return
	}

	// "Unmarshal" the loaded configuration values into our Config struct
	err = viper.Unmarshal(&config)
	return // Return the populated Config struct and any error
}
```

This `LoadConfig` function orchestrates Viper:

1.  It tells Viper where to find the `app.env` file (`AddConfigPath`, `SetConfigName`, `SetConfigType`).
2.  `viper.AutomaticEnv()` enables reading from environment variables (e.g., if `DB_SOURCE` is set in the environment, it might override the value from `app.env`).
3.  `viper.ReadInConfig()` performs the actual reading from the file.
4.  `viper.Unmarshal(&config)` magically maps the values read by Viper into the fields of our `Config` struct.

## Using the Configuration

Once `LoadConfig` successfully runs in `main.go`, we have a `config` variable (of type `util.Config`) holding all our settings. This `config` object is then passed to the different parts of the application that need it.

```go
// --- File: main.go (Simplified startup) ---
package main

import (
	// ... imports
	"github.com/spaghetti-lover/simplebank/util"
	"github.com/spaghetti-lover/simplebank/gapi"
	"github.com/spaghetti-lover/simplebank/worker"
	db "github.com/spaghetti-lover/simplebank/db/sqlc"
	// ...
)

func main() {
	// 1. Load configuration first!
	config, err := util.LoadConfig(".") // "." means look in the current directory
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}

	// 2. Connect to the database using the loaded config
	connPool, err := pgxpool.New(ctx, config.DBSource) // Use config.DBSource
	// ... error handling ...
	store := db.NewStore(connPool)

	// 3. Create other components, passing the config
	redisOpt := asynq.RedisClientOpt{ Addr: config.RedisAddress } // Use config.RedisAddress
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// 4. Pass config to the gRPC server
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)
	// 5. Pass config to the HTTP Gateway server
	runGatewayServer(ctx, waitGroup, config, store, taskDistributor)
	// 6. Pass config to the Task Processor (Worker)
	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)

	// ... wait for servers to run ...
}

// Simplified function signature showing config being passed
func runGrpcServer( /* ... other args ... */, config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	// Create the gRPC server, passing the config object
	server, err := gapi.NewServer(config, store, taskDistributor)
	// ... start server ...
}
```

Now, components like the [gRPC API Service](01_grpc_api_service_.md) can easily access the settings they need from the `Config` struct they receive.

```go
// --- File: gapi/server.go (Simplified) ---
package gapi

import (
	// ... imports
	"github.com/spaghetti-lover/simplebank/token"
	"github.com/spaghetti-lover/simplebank/util"
)

// Server serves gRPC requests.
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config // Holds the application configuration
	store      db.Store
	tokenMaker token.Maker // Our token generator
	// ... other fields
}

// NewServer creates a new gRPC server instance.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	// Use a config value to create the token maker
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey) // Accessing config!
	if err != nil {
		// ... error handling ...
	}

	server := &Server{
		config:     config, // Store the whole config object
		store:      store,
		tokenMaker: tokenMaker,
		// ... initialize other fields ...
	}
	return server, nil
}

// Example usage within an RPC handler (e.g., LoginUser)
func (server *Server) SomeRpcHandler( /* ... */ ) {
	// Access token duration directly from the stored config
	duration := server.config.AccessTokenDuration
	// ... use the duration ...
}
```

The `gapi.Server` stores the `config` object passed to `NewServer`. It can then use values like `config.TokenSymmetricKey` to initialize components (like the `token.Maker`) or `config.AccessTokenDuration` directly within its methods.

## Under the Hood: The Loading Flow

Here's a simplified view of what happens when `main.go` calls `util.LoadConfig`:

```mermaid
sequenceDiagram
    participant Main as main.go
    participant LoadConfig as util.LoadConfig
    participant ViperLib as Viper Library
    participant ConfigSrc as app.env / Env Vars
    participant ConfigStruct as util.Config

    Main->>LoadConfig: Call LoadConfig(".")
    LoadConfig->>ViperLib: Configure paths, name, type ("app.env")
    LoadConfig->>ViperLib: Enable AutomaticEnv()
    LoadConfig->>ViperLib: Call ReadInConfig()
    ViperLib->>ConfigSrc: Read values from app.env file
    ViperLib->>ConfigSrc: Check for overriding Environment Variables
    ConfigSrc-->>ViperLib: Return found values (e.g., DB_SOURCE, TOKEN_SYMMETRIC_KEY)
    ViperLib-->>LoadConfig: Indicate success/failure of reading
    LoadConfig->>ViperLib: Call Unmarshal(&config)
    ViperLib->>ConfigStruct: Populate fields based on loaded values
    ViperLib-->>LoadConfig: Return success/failure of unmarshalling
    LoadConfig-->>Main: Return populated Config struct (or error)
```

1.  `main.go` initiates the process by calling `util.LoadConfig`.
2.  `LoadConfig` configures Viper (where to look, what file name/type, enable env vars).
3.  `LoadConfig` asks Viper to read the configuration (`ReadInConfig`).
4.  Viper reads the `app.env` file and checks for any relevant environment variables.
5.  `LoadConfig` asks Viper to map the loaded values onto the `util.Config` struct (`Unmarshal`).
6.  Viper populates the fields of the `Config` struct.
7.  `LoadConfig` returns the now-populated `Config` struct back to `main.go`.

## Conclusion

Configuration Management is essential for building flexible and secure applications. By separating settings (like database details, secret keys, server addresses) from the code:

- We improve security by not hardcoding sensitive information.
- We make it easy to configure the application for different environments (development, testing, production) without code changes.
- We use **Viper** to read settings from files (`app.env`) and environment variables.
- We load these settings into a Go struct (`util.Config`) for easy and type-safe access throughout the application.

This practice ensures our SimpleBank application knows where its resources are and has the necessary keys to operate, no matter where it's running.

We've now covered almost all the major components of SimpleBank! The final chapter, [Application Orchestration & Entrypoint](08_application_orchestration___entrypoint_.md), will tie everything together, showing how `main.go` acts as the conductor, starting up all the different services (gRPC, HTTP Gateway, Worker) and making sure they run smoothly.

---
