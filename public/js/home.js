(function($) {
    'use strict';
    $(function() {
        var access_token = window.localStorage.getItem('Authorization')

        $('#adminBtn').on("click", function(event) {
            document.location.href='/admin'
              
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
    });
    })(jQuery);