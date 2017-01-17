import { Injectable } from '@angular/core';
import { Http, Response } from '@angular/http';

@Injectable()
export class PushService {
  private serverAPIBaseUrl = 'http://dispotrains.membrives.fr/app/SubscribePush/?endpoint=';

  constructor(private http: Http) {
  }

  registerPushAPI() {
    navigator.serviceWorker.register('service-worker.js', { scope: '/' });
    /*navigator.serviceWorker.register('sw.js', { scope: '/' })
      .then(this.registrationCallback);*/
  }

  private registrationCallback(registration: ServiceWorkerRegistration) {
    registration.pushManager.getSubscription().then((subscription: PushSubscription): Promise<PushSubscription> => {
      if (subscription) {
        return Promise.resolve(subscription);
      }
      return registration.pushManager.subscribe({ userVisibleOnly: false });
    }).then((subscription: PushSubscription) => {
      this.http.get(this.serverAPIBaseUrl + subscription.endpoint);
    });
  }
}
