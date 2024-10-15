import json
from datetime import datetime
from collections import defaultdict

def load_json_file(file_path):
	with open(file_path, 'r') as file:
		data = json.load(file)
	return data

# Convert time from one format to another (returning datetime object)
def convert_time(time_str, from_format="%Y-%m-%dT%H:%M:%SZ", to_format="%m/%d/%Y") -> datetime:
	return datetime.strptime(   (datetime.strptime(time_str, from_format).strftime(to_format))  , to_format   )

json_file_path = "/Users/jameskendrick/Code/cv/features/lineup-generation/v1/static/schedule_raw.json"
data = load_json_file(json_file_path)

schedule = {}

weeks = data["leagueSchedule"]["weeks"]
for week in weeks:
	cur_week = week["weekNumber"]
	start_date = convert_time(week["startDate"])
	end_date = convert_time(week["endDate"])

	# Adjust for all-star break
	if cur_week == 18:
		schedule[17]["endDate"] = end_date
		schedule[17]["gameSpan"] = (end_date - schedule[17]["startDate"]).days
	elif cur_week > 18:
		schedule[cur_week - 1] = {"startDate": start_date,
								  "endDate": end_date,
								  "gameSpan": (end_date - start_date).days}
	else:
		schedule[cur_week] = {"startDate": start_date,
						  	  "endDate": end_date,
						  	  "gameSpan": (end_date - start_date).days}

season_start = datetime.strptime("10/22/2024", "%m/%d/%Y")
game_dates = data["leagueSchedule"]["gameDates"]
cur_week = 1
game_date_format = "%m/%d/%Y %H:%M:%S"

games_in_week = defaultdict(dict)
for day in game_dates:
	game_date = convert_time(day["gameDate"], game_date_format)
	week_start_date = schedule[cur_week]["startDate"]
	week_end_date = schedule[cur_week]["endDate"]
	if game_date < season_start:
		continue
	if game_date > week_end_date:
		week_start_date = schedule[cur_week]["startDate"]
		schedule[cur_week]["games"] = games_in_week
		cur_week += 1
		games_in_week = defaultdict(dict)
	days_since = (game_date - week_start_date).days
	days_since = 0 if days_since == 7 else days_since
	for game in day["games"]:
		games_in_week[game["homeTeam"]["teamTricode"]][int(days_since)] = True
		games_in_week[game["awayTeam"]["teamTricode"]][int(days_since)] = True
# Handle leftover games (ie. last week)
schedule[cur_week]["games"] = games_in_week

# Convert the datetime objects to strings
for week in schedule:
	schedule[week]["startDate"] = schedule[week]["startDate"].strftime("%m/%d/%Y")
	schedule[week]["endDate"] = schedule[week]["endDate"].strftime("%m/%d/%Y")

with open("/Users/jameskendrick/Code/cv/features/lineup-generation/v1/static/schedule.json", 'w') as f:
	json.dump(schedule, f, indent=4)