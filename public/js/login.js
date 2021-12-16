(function($) {
    'use strict';
    $(function() {
        
        $('#loginBtn').on("click", function(event) {
            var inputId = $('#inputId').val();
            var inputPwd = $('#inputPwd').val();
            if(inputId=="" || inputPwd==""){
                event.preventDefault();
                return;
            }
            $.post("/api/login", {id:inputId, password:inputPwd}, getToken);
            })
        var getToken = function(token){
            window.localStorage.setItem("authorization", token.accessUuid)
        }
    
    });
    })(jQuery);