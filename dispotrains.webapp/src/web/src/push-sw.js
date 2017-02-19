self.addEventListener('push', function(event) {
  var notificationOptions = {
    body: "Hello World",
    icon: icon ? icon : 'public/icons/icon-default.png',
    data:{
      url : 'http://example.com/updates'
    }
  };
  title = "Ceci est une notification !";
  return self.registration.showNotification(title, notificationOptions);
});
