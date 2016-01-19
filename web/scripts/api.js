import superagent from 'superagent';

exports.get = function(url, success, failure) {
  superagent
    .get(url)
    .set('Accept', 'application/json')
    .end((err, res) => {
      if (res.ok) {
        success(res);
      } else {
        failure(err, res);
      }
    });
};

exports.put = function(url, data, success, failure) {
  superagent
    .put(url)
    .send(data)
    .set('Accept', 'application/json')
    .end((err, res) => {
      if (res.ok) {
        success(res);
      } else {
        failure(err, res);
      }
    });
};

exports.post = function(url, data, success, failure) {
  superagent
    .post(url)
    .send(data)
    .set('Accept', 'application/json')
    .end((err, res) => {
      if (res.ok) {
        success(res);
      } else {
        failure(err, res);
      }
    });
};

exports.del = function(url, success, failure) {
  superagent
    .del(url)
    .set('Accept', 'application/json')
    .end((err, res) => {
      if (res.ok) {
        success(res);
      } else {
        failure(err, res);
      }
    });
};
