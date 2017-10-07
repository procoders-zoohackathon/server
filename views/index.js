(function () {
    'use strict';
    
    var opts = {}
    opts.extractor = function(obj) {
	return obj["request"]
    };

    var url = 'ws://' + window.location.hostname + ':' + window.location.port + '/connect'
    
    var so = new Socket(url, opts);
    so.On('message', function(data, opts) {
	$('body').append(data+'\n')
    })
})();
