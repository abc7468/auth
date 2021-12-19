(function($) {
    'use strict';
    $(function() {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');

        todoListItem.on('click', '.remove', function() {
            var id = $(this).parent().attr('id')
            var data = {authority: "-1", id:id};

            fetch("/api/users", {
                method: 'DELETE', // or 'PUT'
                body:   JSON.stringify(data)
              }).then($(this).parent().remove())
        });


        todoListItem.on('change', '.checkbox', function() {
            var id = $(this).parent().parent().parent().attr('id')
            console.log(id)

            if ($(this).attr('checked')) {
                var data = {authority: "0", id:id};
                fetch("/api/users", {
                    method: 'UPDATE', // or 'PUT'
                    body:   JSON.stringify(data)
                  }).then($(this).removeAttr('checked'))
                   


                
            } else {
                var data = {authority: "1", id:id};
                fetch("/api/users", {
                    method: 'UPDATE', // or 'PUT'
                    body:   JSON.stringify(data)
                  }).then($(this).attr('checked', 'checked'))
                
            }
        
            
            });



    });
})(jQuery);