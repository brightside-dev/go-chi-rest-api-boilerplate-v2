/*------------------------------------------------------------------
* Bootstrap Simple Admin Template
* Version: 3.0
* Author: Alexis Luna
* Website: https://github.com/alexis-luna/bootstrap-simple-admin-template
-------------------------------------------------------------------*/
(function() {
    'use strict';

    // Toggle sidebar on Menu button click
    $(document).ready(function() {
        $('#sidebarCollapse').on('click', function() {
            $('#sidebar').toggleClass('active');
            $('#body').toggleClass('active');
            console.log('clicked');
        });
    });
    

    // Auto-hide sidebar on window resize if window size is small
    $(window).on('resize', function () {
        if ($(window).width() <= 768) {
            $('#sidebar, #body').addClass('active');
        }
    });

    $(window).on('resize', function () {
        if ($(window).width() > 768) {
            $('#sidebar, #body').removeClass('active');
        }
    });
})();

