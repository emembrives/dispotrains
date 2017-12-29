import { Injectable }     from '@angular/core';
import { Http, Response } from '@angular/http';
import { ElevatorStatistics }  from './elevator-stats';
import { Observable }     from 'rxjs/Observable';
import { combineLatest }  from 'rxjs/observable/combineLatest';
import 'rxjs/add/operator/publishLast';

import { SorterUtils } from './sorting';
import { environment } from '../environments/environment';

@Injectable()
export class ElevatorStatsService {
  private elevatorUrl = environment.baseUrl + '/app/Elevator/';

  constructor(private _http: Http) {}

  getElevatorStats(nameObservable: Observable<string>): Observable<ElevatorStatistics> {
    let self = this;
    return nameObservable.switchMap(function (name) {
      return self._http.get(self.elevatorUrl + name);
    }).catch(this.handleError)
      .map(this.extractData);
  }

  private extractData(res: Response): ElevatorStatistics {
    let body = res.json();
    return new ElevatorStatistics(body);
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
