- Implement the ```/team``` endpoint with the following methods

- **POST**:
	- Create a new team with the given team_name and add the user to the team #inProgress #ishaan
		- Create `team`, with
			- `id = random 6 digit int`
			- `name = team_name`
		- Create relation `teammember`, with:
			- `team_id = team.id`
			- `user_id = user.id`
			- `is_ready = false`
	- Internally assign levels (bosses) to every team on creation #todo
		- Randomly shuffle IDs of `level` 1 to 6
		- Final level will be same for every team
		- Create `teamlevelassignment` relation, with
			- `team_id = team.id`
			- `level_id = level.id`
			- `sequence = level no. for the team`
			- `current_score = 0`
			- `is_finished = 0`
	- Internally assign spells to every team on creation #todo
		- Assign every spell to newly created teams
		- TODO
	- Create websocket channel with `id` as `team.id` and connect every member to it on joining
	
	- Payload: 
		`{"team_name": string}`
	- Response: 
		- success: `200`
		- failure: `{"status":"error", "description": string}`


- **UPDATE**:
	- Find the team by the given team_id and add the current user to the team, provided validity #manan #inProgress 
		- Find `team` where `team.id = id`
		- if valid, create relation `teammember` with:
			- `team_id = team.id`
			- `user_id = user.id`
			- `is_ready = false`
		- Add user to websocket channel with id `team.id` #todo 
	- Validity Checks: 
		- No. of members in team < 4
		- User already not part of any teams
		- Team has not started game
	- Payload: 
		 `{"team_id": int}`
	- Response:
		- success: `200`
		- failure: `{"status":"error", "description": string}`


