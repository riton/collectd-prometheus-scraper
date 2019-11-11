// Copyright (c) IN2P3 Computing Centre, IN2P3, CNRS
// 
// Author(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2019
// 
// This software is governed by the CeCILL-C license under French law and
// abiding by the rules of distribution of free software.  You can  use, 
// modify and/ or redistribute the software under the terms of the CeCILL-C
// license as circulated by CEA, CNRS and INRIA at the following URL
// "http://www.cecill.info". 
// 
// As a counterpart to the access to the source code and  rights to copy,
// modify and redistribute granted by the license, users are provided only
// with a limited warranty  and the software's author,  the holder of the
// economic rights,  and the successive licensors  have only  limited
// liability. 
// 
// In this respect, the user's attention is drawn to the risks associated
// with loading,  using,  modifying and/or developing or reproducing the
// software by the user in light of its specific status of free software,
// that may mean  that it is complicated to manipulate,  and  that  also
// therefore means  that it is reserved for developers  and  experienced
// professionals having in-depth computer knowledge. Users are therefore
// encouraged to load and test the software's suitability as regards their
// requirements in conditions enabling the security of their systems and/or 
// data to be ensured and,  more generally, to use and operate it in the 
// same conditions as regards security. 
// 
// The fact that you are presently reading this means that you have had
// knowledge of the CeCILL-C license and that you accept its terms.

package transport

import (
	"net/http"
	"time"
)

type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

type HTTPBasicCreds struct {
	User     string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type httpClient struct {
	Credentials *HTTPBasicCreds
	client      *http.Client
}

func NewHTTPClient(timeout time.Duration, creds HTTPBasicCreds) *httpClient {
	hClient := http.Client{
		Timeout: timeout,
	}

	return &httpClient{
		Credentials: &creds,
		client:      &hClient,
	}
}

func (hc httpClient) Do(req *http.Request) (*http.Response, error) {
	if hc.Credentials != nil {
		req.SetBasicAuth(hc.Credentials.User, hc.Credentials.Password)
	}
	return hc.client.Do(req)
}