import { Component, Input } from '@angular/core';
import { Observable }               from 'rxjs/Observable';

import { Station } from '../station';

@Component({
  selector: 'station-item',
  templateUrl: './station-item.component.html',
  styleUrls: ['./station-item.component.css'],
})
export class StationItemComponent {
  @Input()
  station: Observable<Station>;

  constructor() {}
}
