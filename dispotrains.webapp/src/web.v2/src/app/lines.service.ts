import { Injectable, Inject } from '@angular/core';
import { Observable, lastValueFrom } from 'rxjs';

import { StationService } from './station.service';
import { Line, Station, NetworkStats } from './station';

class SorterUtils {
  static sorterBySelector(selector: (obj: any) => any) {
    return function(a: any, b: any) {
        if (selector(a) > selector(b)) {
          return 1;
        }

        if (selector(a) < selector(b)) {
          return -1;
        }
        return 0;
      };
  }
}

@Injectable()
export class LinesService {
  constructor(private stationService: StationService) {}

  async getLine(id: string): Promise<Line | undefined> {
    let stations = await this.stationService.getStations()
    let lines = this.stationToLines(stations);
    for (let line of lines) {
      if (line.id === id) {
        return line;
      }
    }
    return undefined;
  }

  getStats(): Promise<NetworkStats> {
    return this.stationService.getStats();
  }

  async getLines(): Promise<Line[]> {
    let stations = await this.stationService.getStations();
    return this.stationToLines(stations);
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

  async getStationsForLine(line: Line): Promise<Station[]> {
    let stations = await this.stationService.getStations();
    return this.findStations(stations, line);
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
