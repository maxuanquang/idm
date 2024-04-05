# IdmProto.IdmServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**idmServiceCreateAccount**](IdmServiceApi.md#idmServiceCreateAccount) | **POST** /api/v1/accounts | 
[**idmServiceCreateDownloadTask**](IdmServiceApi.md#idmServiceCreateDownloadTask) | **POST** /api/v1/tasks | 
[**idmServiceCreateSession**](IdmServiceApi.md#idmServiceCreateSession) | **POST** /api/v1/sessions | 
[**idmServiceDeleteDownloadTask**](IdmServiceApi.md#idmServiceDeleteDownloadTask) | **DELETE** /api/v1/tasks/{downloadTaskId} | 
[**idmServiceDeleteSession**](IdmServiceApi.md#idmServiceDeleteSession) | **DELETE** /api/v1/sessions | 
[**idmServiceGetDownloadTaskFile**](IdmServiceApi.md#idmServiceGetDownloadTaskFile) | **GET** /api/v1/tasks/{downloadTaskId}/files | 
[**idmServiceGetDownloadTaskList**](IdmServiceApi.md#idmServiceGetDownloadTaskList) | **GET** /api/v1/tasks | 
[**idmServiceUpdateDownloadTask**](IdmServiceApi.md#idmServiceUpdateDownloadTask) | **PUT** /api/v1/tasks/{downloadTaskId} | 



## idmServiceCreateAccount

> IdmCreateAccountResponse idmServiceCreateAccount(body)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let body = new IdmProto.IdmCreateAccountRequest(); // IdmCreateAccountRequest | 
apiInstance.idmServiceCreateAccount(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IdmCreateAccountRequest**](IdmCreateAccountRequest.md)|  | 

### Return type

[**IdmCreateAccountResponse**](IdmCreateAccountResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## idmServiceCreateDownloadTask

> IdmCreateDownloadTaskResponse idmServiceCreateDownloadTask(body)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let body = new IdmProto.IdmCreateDownloadTaskRequest(); // IdmCreateDownloadTaskRequest | 
apiInstance.idmServiceCreateDownloadTask(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IdmCreateDownloadTaskRequest**](IdmCreateDownloadTaskRequest.md)|  | 

### Return type

[**IdmCreateDownloadTaskResponse**](IdmCreateDownloadTaskResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## idmServiceCreateSession

> IdmCreateSessionResponse idmServiceCreateSession(body)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let body = new IdmProto.IdmCreateSessionRequest(); // IdmCreateSessionRequest | 
apiInstance.idmServiceCreateSession(body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**IdmCreateSessionRequest**](IdmCreateSessionRequest.md)|  | 

### Return type

[**IdmCreateSessionResponse**](IdmCreateSessionResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json


## idmServiceDeleteDownloadTask

> Object idmServiceDeleteDownloadTask(downloadTaskId)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let downloadTaskId = "downloadTaskId_example"; // String | 
apiInstance.idmServiceDeleteDownloadTask(downloadTaskId, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **downloadTaskId** | **String**|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## idmServiceDeleteSession

> Object idmServiceDeleteSession()



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
apiInstance.idmServiceDeleteSession((error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters

This endpoint does not need any parameter.

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## idmServiceGetDownloadTaskFile

> StreamResultOfIdmGetDownloadTaskFileResponse idmServiceGetDownloadTaskFile(downloadTaskId)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let downloadTaskId = "downloadTaskId_example"; // String | 
apiInstance.idmServiceGetDownloadTaskFile(downloadTaskId, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **downloadTaskId** | **String**|  | 

### Return type

[**StreamResultOfIdmGetDownloadTaskFileResponse**](StreamResultOfIdmGetDownloadTaskFileResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## idmServiceGetDownloadTaskList

> IdmGetDownloadTaskListResponse idmServiceGetDownloadTaskList(opts)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let opts = {
  'offset': "offset_example", // String | 
  'limit': "limit_example" // String | 
};
apiInstance.idmServiceGetDownloadTaskList(opts, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **offset** | **String**|  | [optional] 
 **limit** | **String**|  | [optional] 

### Return type

[**IdmGetDownloadTaskListResponse**](IdmGetDownloadTaskListResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json


## idmServiceUpdateDownloadTask

> IdmUpdateDownloadTaskResponse idmServiceUpdateDownloadTask(downloadTaskId, body)



### Example

```javascript
import IdmProto from 'idm_proto';

let apiInstance = new IdmProto.IdmServiceApi();
let downloadTaskId = "downloadTaskId_example"; // String | 
let body = new IdmProto.IdmServiceUpdateDownloadTaskBody(); // IdmServiceUpdateDownloadTaskBody | 
apiInstance.idmServiceUpdateDownloadTask(downloadTaskId, body, (error, data, response) => {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
});
```

### Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **downloadTaskId** | **String**|  | 
 **body** | [**IdmServiceUpdateDownloadTaskBody**](IdmServiceUpdateDownloadTaskBody.md)|  | 

### Return type

[**IdmUpdateDownloadTaskResponse**](IdmUpdateDownloadTaskResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

