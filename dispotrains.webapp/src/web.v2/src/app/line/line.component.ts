import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { LinesService } from '../lines.service';
import { Line, Station } from '../station';

@Component({
  selector: 'app-line',
  templateUrl: './line.component.html',
  styleUrls: ['./line.component.css']
})
export class LineComponent implements OnInit {
  line: Line | undefined;
  badStations: Station[] | undefined;
  goodStations: Station[] | undefined;

  constructor(
    private linesService: LinesService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.params.subscribe(async (params) => {  
      // Do Something with the params you receive
      let lineId = params['id'];
      this.line = await this.linesService.getLine(lineId)
      if (this.line === undefined) { return; }
      let stations = await this.linesService.getStationsForLine(this.line);
      this.goodStations = this.findGoodStations(stations);
      this.badStations = this.findBadStations(stations);
    })
  }

  private findGoodStations(stations: Station[]) : Station[] {
    return stations.filter((station: Station) => station.available());
  }

  private findBadStations(stations: Station[]) : Station[] {
    return stations.filter((station: Station) => !station.available());
  }
}
