import { Injectable, Inject } from '@angular/core';
import { Observable } from 'rxjs/Observable';
import { combineLatest } from 'rxjs/observable/combineLatest';

import { StationService } from './station.service';
import { Line, Station } from './station';
import { SorterUtils } from './sorting';

@Injectable()
export class LinesService {
  constructor(private stationService: StationService) {}

  getLine(id: string): Promise<Line> {
    return this.stationService.getStations()
      .then(this.stationToLines)
      .then(this.findLine(id));
  }

  private findLine(id: string): (lines: Line[]) => Line {
    return function(lines: Line[]) {
      for (let line of lines) {
        if (line.id === id) {
          return line;
        }
      }
    };
  }

  getLines(): Promise<Line[]> {
    return this.stationService.getStations().then(this.stationToLines);
  }

  private stationToLines(stations: Station[]): Line[] {
    let lineMap: Map<string, Line> = new Map<string, Line>();
    for (let station of stations) {
        for (let line of station.lines) {
          let lineName: string = line.GetName();
          lineMap.set(lineName, line);
        }
    }
    let lines: Line[] = new Array<Line>();
    lineMap.forEach(function(value: Line, key: string) { lines.push(value); });
    lines.sort(SorterUtils.sorterBySelector((line: Line) => line.GetName()));
    return lines;
  }

  getStationsForLine(line: Observable<Line>): Observable<Station[]> {
    return combineLatest(this.stationService.getStations(), line, this.findStations);
  }

  private findStations(stations: Station[], line: Line) : Station[] {
    let selected: Station[] = new Array<Station>();
    for (let station of stations) {
      if (station.lines.some((stationLine) => stationLine.id == line.id)) {
        selected.push(station);
      }
    }
    selected.sort(SorterUtils.sorterBySelector((station: Station) => station.displayname));
    return selected;
  }
}
