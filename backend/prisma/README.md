# Labyrinth 2025 Database Schema


# Conventions to be followed
1. Every model definition must be in PascalCase
2. Every model must be mapped to all lower case (Example: TeamMemberAssignment -> teammemberassignment)
3. Every attribute that can be mapped must be mapped to snake_case
4. All @relations must be defined after all attributes, a // must separate attributes and relations
5. `npx prisma format` must be run before every commit


# Documentation

## Spell / Progress flow:

1. Teams are assignment spells at first (Entries made in TeamSpellAssignment)
2. Team picks spell to cast (Entry made in TeamSpellAttempt as `active=true`)
   - If team unselects, previous entry is updated to `active=false`
3. If them proceeds, they go to locations
4. On location, when qr code is scanned, entry made in TeamLocationProgress tied to the attempt in which the team scanned code.
5. When `n` locations scanned, we take locations, pull all patterns corresponding to Spell from SpellPattern, and check if our locations match up to any pattern.


# Notes
1. Location \<m-----n> Pattern \<m-----1> Spell