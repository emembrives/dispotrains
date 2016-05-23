import json
import datetime
import itertools

statuses = []
f = open('statuses-2014.json')
for line in f.readlines():
    statuses.append(json.loads(line.strip()))
f.close()

class StatusByDate(object):
    def __init__(self):
        self.sncf = 0
        self.ratp = 0

def filter_by_date(start_date, end_date):
    start = datetime.datetime.strptime(start_date, "%Y-%m-%d")
    end = datetime.datetime.strptime(end_date, "%Y-%m-%d")

    def filter_by_provided_date(entry):
        entry_date = datetime.datetime.fromtimestamp(int(entry["lastupdate"]["$date"])/1000)
        return start <= entry_date and entry_date <= end
    return filter_by_provided_date

def to_simplified_status(entry):
    new_entry = {"elevator": entry["elevator"], "state": entry["state"]}
    new_entry["date"] = datetime.datetime.fromtimestamp(int(entry["lastupdate"]["$date"])/1000)
    return new_entry

filtered_statuses = itertools.ifilter(filter_by_date("2014-01-01", "2014-07-01"), statuses)
simplified_statuses = itertools.imap(to_simplified_status, filtered_statuses)

statuses_by_date = {}
for status in simplified_statuses:
    status_date = datetime.date(status["date"].year, status["date"].month, status["date"].day)
    date_entry = statuses_by_date.setdefault(status_date, StatusByDate())
    if status["elevator"].isdigit():
        date_entry.sncf += 1
    else:
        date_entry.ratp += 1

for k, v in statuses_by_date.items():
    print k.isoformat(), v.sncf, v.ratp, v.sncf+v.ratp
