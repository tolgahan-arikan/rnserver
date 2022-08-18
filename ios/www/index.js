document.getElementById('newStuff').innerHTML = 'The javascript works?';
fetch('https://reqres.in/api/products/3')
  .then(response => response.json())
  .then(data => {
    document.getElementById('newStuff2').innerHTML =
      'yep! <br/>' + JSON.stringify(data, null, 2);
  })
  .catch(error => {
    document.getElementById('newStuff').innerHTML = error;
  });
