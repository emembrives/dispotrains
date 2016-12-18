// Add the RxJS Observable operators.
import './rxjs-operators';

import { Component, OnInit } from '@angular/core';
import { Observable }     from 'rxjs/Observable';

import { StationService } from './station.service';
import { Station } from './station';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
  providers: [StationService]
})
export class AppComponent implements OnInit {
  stations: Observable<Station[]>;

  title = 'app works!';

  constructor(private stationService: StationService) { }

  ngOnInit(): void {
    this.stations = this.stationService.getStations();
  }

  gotoDetail(station: Station): void {
    // Route to detail
  }
}
