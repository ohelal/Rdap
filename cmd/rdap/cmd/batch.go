package cmd

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strings"
    "sync"
    "github.com/spf13/cobra"
)

var batchCmd = &cobra.Command{
    Use:   "batch [file]",
    Short: "Batch query from file",
    Long: `Query RDAP information for multiple items from a file.
Each line in the file should be a domain name, IP address, or ASN.
Example: rdap batch queries.txt`,
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
        filePath := args[0]

        // Read queries from file
        file, err := os.Open(filePath)
        if err != nil {
            return fmt.Errorf("opening file: %w", err)
        }
        defer file.Close()

        var queries []string
        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            if line := strings.TrimSpace(scanner.Text()); line != "" {
                queries = append(queries, line)
            }
        }

        if err := scanner.Err(); err != nil {
            return fmt.Errorf("reading file: %w", err)
        }

        // Process queries concurrently
        results := make(map[string]interface{})
        var wg sync.WaitGroup
        var mu sync.Mutex

        if verbose {
            fmt.Printf("Processing %d queries...\n", len(queries))
        }

        for _, query := range queries {
            wg.Add(1)
            go func(q string) {
                defer wg.Done()

                var result interface{}
                var err error

                ctx := context.Background()
                switch {
                case strings.HasPrefix(strings.ToUpper(q), "AS"):
                    result, err = client.QueryASN(ctx, q)
                case strings.Contains(q, "."):
                    result, err = client.QueryDomain(ctx, q)
                default:
                    result, err = client.QueryIP(ctx, q)
                }

                if err != nil {
                    result = fmt.Sprintf("Error: %v", err)
                }

                mu.Lock()
                results[q] = result
                mu.Unlock()
            }(query)
        }

        wg.Wait()

        // Render results
        fmt.Printf("\n=== Batch Query Results ===\n")
        for query, result := range results {
            fmt.Printf("\n[Query: %s]\n", query)
            if data, ok := result.(map[string]interface{}); ok {
                fmt.Print(formatRDAPResult(data))
            } else {
                fmt.Printf("%v\n", result)
            }
            fmt.Println(strings.Repeat("-", 80))
        }

        return nil
    },
}

func init() {
    rootCmd.AddCommand(batchCmd)
}
