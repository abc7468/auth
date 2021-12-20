(function($) {
    'use strict';
    $(function() {
        var access_token = window.localStorage.getItem('Authorization')

        $('#adminBtn').on("click", function(event) {
          fetch("/auth/user", {
            method: 'GET', // or 'PUT'
            headers:{
                "Authorization":access_token
            }
          }).then(res=>res.json())
          .then(function(response){
            if(response.user_authorized=="1"){
              alert("권한 부족")
              return
            }
            else{
              console.log(response.user_authorized)
              document.location.href = "/admin"

            }
          })
        });
        $('#logoutBtn').on("click", function(event) {
        
            fetch("/auth/token", {
                method: 'DELETE', // or 'PUT'
                headers:{
                    "Authorization":access_token
                }
              }).then(res=>res.json())
                .then( window.localStorage.removeItem('Authorization'))
              .then(
                setTimeout(function() {
                  document.location.href = "/"
               }, 500))

        })


        $(document).ready(function(){

          fetch('/auth/token/valid', {
              method: "GET",
              headers: {
                "Authorization": access_token
              }
            }).then(res => res.json())
            .then(function(response){
              if(response.access_token!=""){
                console.log(response.access_token)
                window.localStorage.setItem("Authorization", "bearer " + response.access_token)
                return
              }  
              if(!response.success){
                console.log("false")
                  window.localStorage.removeItem('Authorization')
                  window.location.href ='/'
                  return
              }
                
            })
      })   
    });
    })(jQuery);