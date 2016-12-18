import { Injectable } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { StationService } from './station.service';
import { Line, Station } from './station';

@Injectable()
export class LinesService {
  constructor(private stationService: StationService) { }

  getLines() : Observable<Line[]> {
    return this.stationService.getStations().map(this.stationToLines);
  }

  private stationToLines(stations: Station[]) : Line[] {
    let lineSet: Set<Line>;
    for (let station of stations) {
        for (let line of station.lines) {
          lineSet.add(line);
        }
    }
    let lines: Line[];
    lineSet.forEach(function(line: Line) { lines.push(line); });
    return lines;
  }

  getStationsForLine(line: Line) : Observable<Station[]> {
    return this.stationService.getStations().map(this.findStations(line));
  }

  private findStations(line: Line) : ((stations: Station[])=>Station[]) {
    return function(stations: Station[]) : Station[] {
      let selected: Station[];
      for (let station of stations) {
        if (station.lines.some((stationLine) => stationLine == line)) {
          selected.push(station);
        }
      }
      return selected;
    }
  }
}
