# Project Instructions for Gemini

## Technology Stack

### Frontend
- Next.js
- TypeScript

### Backend
- Go (Gin)

---

## Architecture
- Domain-Driven Design (DDD)
- AWS Lambda

---

## Instructions for Code Changes

- After implementing changes, always explain what was changed in **Japanese**.
- All comments in the implementation must be written in **Japanese**.
- Do not add comments such as "Added", "New", or similar for newly added functions or parameters.
- Do not add comments such as "Removed", "Deleted", or similar for removed functions or parameters.
- Keep comments to a minimum and only describe what the function does.

---

## Critical Rule: back/.env

- **Never edit, rewrite, delete, or normalize `back/.env`.**
- This file contains sensitive environment variables and must not be touched.
- Assume `back/.env` already exists and is correctly configured.
- If environment variables are required, reference them conceptually without modifying the file.
