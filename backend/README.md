# To Build

```bash
  go build
```

# To run
```

./labyrinth

```

# Supabase Local Development
To push all tables and triggers to the supabase instance:
```bash
npx supabase start
npx supabase db reset
npx prisma db push
npx prisma db seed
```