import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import { Station } from '../station';
import { StationService } from '../station.service';

@Component({
  selector: 'app-station',
  templateUrl: './station.component.html',
  styleUrls: ['./station.component.css']
})
export class StationComponent implements OnInit {
  station: Station | undefined;
  lineId: string | undefined

  constructor(
    private stationService: StationService,
    private route: ActivatedRoute,
  ) { }

  ngOnInit() {
    this.route.params
      .subscribe(async (params: Params) => {
        this.lineId = params['lineId'];
        let stationId: string = params['id'];
        this.station = await this.stationService.getStation(stationId);
      });
  }
}
