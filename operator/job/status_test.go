package job

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	k8upv1 "github.com/k8up-io/k8up/v2/api/v1"
)

func newStatusTestConfig(t *testing.T) (Config, *k8upv1.Backup) {
	scheme := runtime.NewScheme()
	require.NoError(t, k8upv1.AddToScheme(scheme))
	backup := &k8upv1.Backup{ObjectMeta: metav1.ObjectMeta{Name: "backup", Namespace: "ns"}}
	fclient := fake.NewClientBuilder().
		WithScheme(scheme).
		WithStatusSubresource(&k8upv1.Backup{}).
		WithObjects(backup).
		Build()
	return NewConfig(fclient, backup, ""), backup
}

func TestSetStarted_PersistsStatus(t *testing.T) {
	c, backup := newStatusTestConfig(t)

	c.SetStarted(context.Background(), "started %s", "now")

	stored := &k8upv1.Backup{}
	require.NoError(t, c.Client.Get(context.Background(), client.ObjectKeyFromObject(backup), stored))
	assert.True(t, stored.Status.Started)
	assert.True(t, stored.Status.HasStarted())
}

func TestSetFinished_PersistsStatus(t *testing.T) {
	c, backup := newStatusTestConfig(t)

	c.SetFinished(context.Background(), "ns", "backup")

	stored := &k8upv1.Backup{}
	require.NoError(t, c.Client.Get(context.Background(), client.ObjectKeyFromObject(backup), stored))
	assert.True(t, stored.Status.Finished)
}
