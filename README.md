# Labyrinth 2025 Database Schema

# Conventions to be followed
1. Every model definition must be in PascalCase
2. Every model must be mapped to all lower case (Example: TeamMemberAssignment -> teammemberassignment)
3. Every attribute that can be mapped must be mapped to snake_case
4. All @relations must be defined after all attributes, a // must separate attributes and relations
5. `npx prisma format` must be run before every commit