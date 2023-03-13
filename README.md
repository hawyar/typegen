### typepgen

> Generate typescript declaration file from JSON Schema

Given: `user.schema.json`

```json
{
  "title": "User",
  "type": "object",
  "properties": {
    "id": {
      "type": "string"
    },
    "name": {
      "type": "string"
    },
    "age": {
      "type": "number"
    }
  },
  "required": ["id"]
}
```
Run typegen

```bash
typegen user.schema.json
```

Outputs declaration file: `User.d.ts`

```typescript
export interface User {
  name?: string;
  age?: number;
  id: string;
}

## Build
  
```bash
make build
```

```



