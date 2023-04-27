/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhook

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func RegistryHttpServer(port uint16) error {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	r.POST("/webhook", webhookHandler)
	// k8s core group
	apiCore := r.Group("/api/v1")
	{
		apiCore.GET("/namespaces/:namespace/pods/:podName/log", podLogs)
	}
	return r.Run(fmt.Sprintf("%s:%d", "0.0.0.0", port))
}