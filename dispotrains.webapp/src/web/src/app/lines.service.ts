import { Injectable, Inject } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { StationService } from './station.service';
import { Line, Station } from './station';

@Injectable()
export class LinesService {
  constructor(@Inject(StationService) private stationService: StationService) {}

  getLines() : Observable<Line[]> {
    return this.stationService.getStations().map(this.stationToLines);
  }

  private stationToLines(stations: Station[]) : Line[] {
    let lineMap: Map<string, Line> = new Map<string, Line>();
    for (let station of stations) {
        for (let line of station.lines) {
          let lineName: string = line.GetName();
          lineMap.set(lineName, line);
        }
    }
    let lines: Line[] = new Array<Line>();
    lineMap.forEach(function(value: Line, key: string) { lines.push(value); });
    lines.sort((a: Line, b: Line) => {
      if (a.GetName() > b.GetName()) {
        return 1;
      }

      if (a.GetName() < b.GetName()) {
        return -1;
      }
      return 0;
    });
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
