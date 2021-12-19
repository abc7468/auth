(function($) {
    'use strict';
    $(function() {
        
        $('#signUpBtn').on("click", function(event) {
            var inputEmail = $('#email').val();
            var inputPwd = $('#pwd').val();
            var inputCheckPwd = $('#pwd_check').val();
            var inputName = $('#name').val();
            var inputPhone = $('#phone').val();
            if(inputCheckPwd!=inputPwd){
                alert("값 다름")
                return
            }
            $.post("/api/signup", {email:inputEmail, password:inputPwd, name:inputName, phone:inputPhone}, test);
            })
        
            var test = function(res){
                if(res.error){
                    alert(res.error)
                    return
                }
                window.location.href = '/'
            }
    
    });
    })(jQuery);