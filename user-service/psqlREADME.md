## 🔑 Goal

You want to open a terminal where you can enter SQL like:

```sql
CREATE USER user_service WITH PASSWORD 'password';
```

But **how you get into that terminal (`psql`) depends on your operating system and setup**.

---

## 🧑💻 macOS (Homebrew)

### ✅ Step-by-Step

1. **Start PostgreSQL** if it's not already running:

   ```bash
   brew services start postgresql
   ```

2. **Open your Terminal** (e.g., Terminal.app, iTerm).

3. **Enter the `psql` shell** by running:

   ```bash
   psql postgres
   ```

   * This logs you into the `postgres` database using your macOS username.

4. Now you're inside the PostgreSQL shell and can run:

   ```sql
   CREATE USER user_service WITH PASSWORD 'password';
   CREATE DATABASE user_service_db OWNER user_service;
   ```

5. Type `\q` to quit the shell.

---

## 🪟 Windows

### ✅ Step-by-Step

1. After installing PostgreSQL using the installer, open:

   ```
   Start Menu → PostgreSQL → SQL Shell (psql)
   ```

2. A prompt will ask you for:

   ```
   Server [localhost]:
   Database [postgres]:
   Port [5432]:
   Username [postgres]:
   Password:
   ```

   * Just press Enter for default values unless you changed them.
   * Enter the **password** you set during installation when prompted.

3. You are now in the `psql` shell and can run:

   ```sql
   CREATE USER user_service WITH PASSWORD 'password';
   CREATE DATABASE user_service_db OWNER user_service;
   ```

---

## 🐧 Linux (Ubuntu/Debian)

### ✅ Step-by-Step

1. Start PostgreSQL if it's not already running:

   ```bash
   sudo service postgresql start
   ```

2. Switch to the PostgreSQL system user:

   ```bash
   sudo -i -u postgres
   ```

3. Enter the `psql` shell:

   ```bash
   psql
   ```

4. You’re now in the shell and can run:

   ```sql
   CREATE USER user_service WITH PASSWORD 'password';
   CREATE DATABASE user_service_db OWNER user_service;
   ```

5. Type `\q` to exit.

---

## 🧪 Test Connection

After creating the user and database, verify it works:

```bash
psql -U user_service -d user_service_db
```

If you see the `psql` prompt without errors, it worked!

---
