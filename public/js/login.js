(function($) {
    'use strict';
    $(function() {
      var access_token = window.localStorage.getItem('Authorization')
        
        $('#loginBtn').on("click", function(event) {
            var inputEmail = $('#inputEmail').val();
            var inputPwd = $('#inputPwd').val();
            var data = {email: inputEmail, password:inputPwd};
            if(inputEmail=="" || inputPwd==""){
                event.preventDefault();
                return;
            }
            fetch("/api/login", {
                method: 'POST', // or 'PUT'
                body: JSON.stringify(data), // data can be `string` or {object}!
                headers:{
                  'Content-Type': 'application/json'
                }
              }).then(res=>res.json())
              .then(function(response){
                if(response.access_token){
                  window.localStorage.setItem("Authorization", "bearer " + response.access_token)
                }
              }).then(
                setTimeout(function() {
                  window.location.href = "/main"
               }, 500)
              ).catch(error => console.error('Error:', error))

    
        });

        $(document).ready(function(){

            fetch('/auth/token', {
                method: "GET",
                headers: {
                  "Authorization": access_token
                }
              }).then(res => res.json())
              .then(function(response){
                  if(response.success){
                    window.location.href ='/main'
                  }
              })
        })    
    })
    })(jQuery);

