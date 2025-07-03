# NStore

NStore is a lightweight storage module that lets you persist structured data as JSON files on disk while keeping them loaded in memory for fast and reusable access. It's ideal for managing dynamic configurations, keys, scripts, or any data that needs simple persistence without the overhead of a traditional database.

Unlike a database, NStore focuses on simplicity: all data is saved as plain JSON files and automatically loaded into memory, making it globally accessible across your application. This design supports a singleton-like pattern for centralized access without external dependencies.

### Usage

#### Define your structure to storage
```go
type Keys struct {
	*nstore.Metadata //mandatory metadata for storage management

	Name     string `json:"name"`
	KeyData  *Data  `json:"data"`
}
```

#### Instance of NStorage

```go
folder := "managed/keys"
storage, err := nstore.New[*Keys](folder)

storage.LoadFromDisk()

err := storage.Save(&Keys{Name: "my-key", KeyData: &Data{...} })

storage.Exists("<ID>")

key, err := storage.Load("<ID>")

err := storage.Delete("<ID>")

storage.ListOfCache()

storage.IDs() // ["4ba30010-5f83-46f0-ab7e-6cbb566e1f0d", ...]

results, length := storage.Query(func(t *User) bool { return t.Name == "my-key" }, 3)
```

### Example JSON file
```json
{
  "id": "4ba30010-5f83-46f0-ab7e-6cbb566e1f0d",
  "rev": "514eb39c-f5d1-40d6-b132-7b60cf68b87b",
  "created_at": "2025-07-03T00:44:14.33842-03:00",
  "modified_at": "2025-07-03T00:44:14.33842-03:00",
  "name": "my-key",
  "data": {
    ...
  }
}
```