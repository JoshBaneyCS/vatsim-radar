import (
    "github.com/wailsapp/wails/v2/pkg/runtime"
    // ...other imports
)

func main() {
    app := NewApp() // our App struct for Wails (if needed)
    // Set up Wails application options (including Title, window size, etc, see Branding section below)
    err := wails.Run(&options.App{
        Title:  appTitle,
        Width:  1280,
        Height: 720,
        // ... other options like AssetServer, OnStartup, Bind (if needed)
        OnStartup: func(ctx context.Context) {
            // Start data update loop when app starts
            go func() {
                for {
                    // fetch VATSIM pilots
                    pilots, err1 := fetchVatsimData()
                    if err1 == nil {
                        // Emit event "vatsim-data" with the pilot list
                        runtime.EventsEmit(ctx, "vatsim-data", pilots)
                    }
                    // fetch Real flights
                    flights, err2 := fetchRealTraffic()
                    if err2 == nil {
                        runtime.EventsEmit(ctx, "real-data", flights)
                    }
                    // Sleep until next update
                    time.Sleep(15 * time.Second)
                }
            }()
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
