/*
Copyright 2018 Pressinfra SRL

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

package wordpress

import (
	"reflect"

	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"

	wpapi "github.com/presslabs/wordpress-operator/pkg/apis/wordpress/v1alpha1"
)

func (c *Controller) initWordpressWorker() {
	c.wpInformer = c.wpInformerFactory.Wordpress().V1alpha1().Wordpresses().Informer()
	c.wpLister = c.wpInformerFactory.Wordpress().V1alpha1().Wordpresses().Lister()
	c.wpQueue = queue.New("Wordpress", maxRetries, threadiness, c.reconcileWordpress)
	c.wpInformer.AddEventHandler(queue.NewEventHandler(c.wpQueue.GetQueue(), func(old interface{}, new interface{}) bool {
		oldSpec, ok := old.(*wpapi.Wordpress)
		if !ok {
			return false
		}
		newSpec, ok := new.(*wpapi.Wordpress)
		if !ok {
			return false
		}
		return !reflect.DeepEqual(oldSpec.Spec, newSpec.Spec)
	}))
}

func (c *Controller) reconcileWordpress(key string) error {
	obj, exists, err := c.wpInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}
	if exists {
		glog.Infof("Sync/Add/Update for Wordpress %s", key)
		wp := obj.(*wpapi.Wordpress).DeepCopy()

		c.syncDeployment(wp)
	}
	return nil
}