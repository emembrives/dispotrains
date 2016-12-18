import { Component, OnInit, Input } from '@angular/core';
import { Station } from '../station';

@Component({
  selector: 'station-item',
  templateUrl: './station-item.component.html',
  styleUrls: ['./station-item.component.css']
})
export class StationItemComponent implements OnInit {
  @Input()
  station: Station;

  constructor() { }

  ngOnInit() {
  }

}
