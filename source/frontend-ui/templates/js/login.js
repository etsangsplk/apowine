function signOut() {
  var auth2 = gapi.auth2.getAuthInstance();
  auth2.signOut().then(function () {
    alert("User Successfully Signed Out")
  });
  var xhr = new XMLHttpRequest();
  xhr.open('POST', '/catchtoken');
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.send('authenticated=' + false);
}

function onLoad() {
  gapi.load('auth2', function() {
    gapi.auth2.init();
  });
}
