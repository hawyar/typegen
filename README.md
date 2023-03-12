Generate TS declaration files `.d.ts` from JSON Schema


Input: user.schema.json
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

```bash
typegen user.schema.json
```

Output:
```typescript
export interface User {
	name?: string;
	age?: number;
	id: string;
}
```



