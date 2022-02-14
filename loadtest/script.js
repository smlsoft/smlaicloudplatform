import http from 'k6/http';

export default function (data) {

  const url = 'http://localhost:8089/merchant/23xK48ZSaDPzoxZVXIbV8w6kFVw/inventory';
  const payload = JSON.stringify(
    {
			"itemSku":      "devx001",
			"merchantId":   "23twO9nFtgsLGAuQ9JXPzi3C65N",
			"categoryGuid": "cat1234",
			"lineNumber":   1,
			"price":        99.0,
			"recommended":  false,
			"haveImage":    false,
			"activated":    true,
			"name1":        "Name 1",
			"description1": "Description 1",
    }
  );

  const params = {
      headers: {
          'Content-Type': 'application/json',
          "Authorization": "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NDQ2MjkxMzMsInVzZXJuYW1lIjoiZGV2MDEiLCJuYW1lIjoiZGV2IGRldiJ9.f4dIsHbKaH7K3XrUa39Ft7-m8kpZoIC_CwykDktu9KCB0AbY4bwATYJyaeAGokHa_Fy70jLI4WedCqM7XGZ5w20HkjZm4ZsA-XYRQDCK4WkJ0Elc6L6yEBDk2kr2oqjfLQV3WL4dbUBa_bmkIxU0D0GAH3hjMQF8hVbazHtrJ0t0f86IFqvqnSNoBmmWJy-w43hFDDf7G8v5L_vK-nUKI9wY0t1bdZ_-2fVfWd_0An9dCcgQ4qIce9ad1jf4P6muI1YiZWX4bHjYCIURUVB8Ch4rPW2ekeqOvN1HWxFIFkjdtOM7gAlTsmnmKeSArGB_7RIiNqL-19lKi8gyqFhBe-dcZHOUBraidpjH7ReMr-w9AAY6fXAE8GyXsDV--0eGho1I9UdDqYQa4zonVMUW1xncS3KWwVxrdfRrW-KIvWc6JBXgpAdH_jeLf-vXywHfArAXkLCJndYSCY4RFgOxPO39w3NoYTBLuMB0liIWSwJbdjxXsBPTnJVBWRsO-U_lnlN5IdZyP22JVR0pkXINKQDSS5agobFHnavZp-s2ckjuqquSwbn4tkbJwDGDVybVKsBycdvt5Tminb0KlD8YWOjF1HfKAY0iU84_RDlJ9L1tJXhD_4k1_Bd-9JfymDexHKkdOpKjKhLWHnD5SozIrPwS_CjBjCylr4rv_5NwEWs",
      }
  };

  http.post(url, payload, params);
}
