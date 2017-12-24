import { Injectable }     from '@angular/core';
import { Http, Response } from '@angular/http';
import { ElevatorStats }  from './elevator-stats';
import { Observable }     from 'rxjs/Observable';
import { combineLatest }  from 'rxjs/observable/combineLatest';
import 'rxjs/add/operator/publishLast';

import { SorterUtils } from './sorting';
import { environment } from '../environments/environment';

@Injectable()
export class ElevatorStatsService {
  private elevatorUrl = environment.baseUrl + '/app/GetElevator/';

  constructor(private http: Http) {}

  getElevatorStats(nameObservable: Observable<string>): Observable<ElevatorStats[]> {
    return nameObservable.map(
        function(name: String) {return this.http.get(this.elevatorUrl + name); })
      .catch(this.handleError)
      .map(this.extractData)
      .publishLast()
      .refCount();
  }

  private extractData(res: Response): ElevatorStats[] {
    let body = res.json();
    let elevatorStats = new Array<ElevatorStats>();
    for (let elevData of body) {
      elevatorStats.push(new ElevatorStats(elevData));
    }
    return elevatorStats;
  }

  private handleError(error: Response | any) {
    // In a real world app, we might use a remote logging infrastructure
    let errMsg: string;
    if (error instanceof Response) {
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error('Error while retrieving stations: ' + errMsg);
    return Observable.throw(errMsg);
  }
}
