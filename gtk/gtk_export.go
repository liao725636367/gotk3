package gtk

/*
 #include <gtk/gtk.h>
*/
import "C"
import (
	"strings"
	"unsafe"

	"github.com/gotk3/gotk3/glib"
)

//export substring_match_equal_func
func substring_match_equal_func(model *C.GtkTreeModel,
	column C.gint,
	key *C.gchar,
	iter *C.GtkTreeIter,
	data C.gpointer) C.gboolean {

	goModel := &TreeModel{glib.Take(unsafe.Pointer(model))}
	goIter := &TreeIter{(C.GtkTreeIter)(*iter)}

	value, err := goModel.GetValue(goIter, int(column))
	if err != nil {
		return gbool(true)
	}

	str, _ := value.GetString()
	if str == "" {
		return gbool(true)
	}

	subStr := C.GoString((*C.char)(key))
	res := strings.Contains(str, subStr)
	return gbool(!res)
}

//export goBuilderConnect
func goBuilderConnect(builder *C.GtkBuilder,
	object *C.GObject,
	signal_name *C.gchar,
	handler_name *C.gchar,
	connect_object *C.GObject,
	flags C.GConnectFlags,
	user_data C.gpointer) {

	builderSignals.Lock()
	signals, ok := builderSignals.m[builder]
	builderSignals.Unlock()

	if !ok {
		panic("no signal mapping defined for this GtkBuilder")
	}

	h := C.GoString((*C.char)(handler_name))
	s := C.GoString((*C.char)(signal_name))

	handler, ok := signals[h]
	if !ok {
		return
	}

	if object == nil {
		panic("unexpected nil object from builder")
	}

	//TODO: figure out a better way to get a glib.Object from a *C.GObject
	gobj := glib.Object{glib.ToGObject(unsafe.Pointer(object))}
	gobj.Connect(s, handler)
}

//export goPageSetupDone
func goPageSetupDone(setup *C.GtkPageSetup,
	data C.gpointer) {

	id := int(uintptr(data))

	pageSetupDoneCallbackRegistry.Lock()
	r := pageSetupDoneCallbackRegistry.m[id]
	delete(pageSetupDoneCallbackRegistry.m, id)
	pageSetupDoneCallbackRegistry.Unlock()

	obj := glib.Take(unsafe.Pointer(setup))
	r.fn(wrapPageSetup(obj), r.data)

}

//export goPrintSettings
func goPrintSettings(key *C.gchar,
	value *C.gchar,
	userData C.gpointer) {

	id := int(uintptr(userData))

	printSettingsCallbackRegistry.Lock()
	r := printSettingsCallbackRegistry.m[id]
	// TODO: figure out a way to determine when we can clean up
	//delete(printSettingsCallbackRegistry.m, id)
	printSettingsCallbackRegistry.Unlock()

	r.fn(C.GoString((*C.char)(key)), C.GoString((*C.char)(value)), r.userData)

}

//export goTreeModelFilterFuncs
func goTreeModelFilterFuncs(filter *C.GtkTreeModelFilter, iter *C.GtkTreeIter, data C.gpointer) C.gboolean {
	id := int(uintptr(data))

	treeModelVisibleFilterFuncRegistry.Lock()
	r := treeModelVisibleFilterFuncRegistry.m[id]
	treeModelVisibleFilterFuncRegistry.Unlock()

	goIter := &TreeIter{(C.GtkTreeIter)(*iter)}
	return gbool(r.fn(
		wrapTreeModelFilter(glib.Take(unsafe.Pointer(filter))),
		goIter,
		r.userData))
}

//export goTreeSortableSortFuncs
func goTreeSortableSortFuncs(model *C.GtkTreeModel, a, b *C.GtkTreeIter, data C.gpointer) C.gint {
	id := int(uintptr(data))

	treeStoreSortFuncRegistry.Lock()
	r := treeStoreSortFuncRegistry.m[id]
	treeStoreSortFuncRegistry.Unlock()

	goIterA := &TreeIter{(C.GtkTreeIter)(*a)}
	goIterB := &TreeIter{(C.GtkTreeIter)(*b)}

	return C.gint(r.fn(wrapTreeModel(glib.Take(unsafe.Pointer(model))), goIterA, goIterB, r.userData))
}

//export goTreeModelForeachFunc
func goTreeModelForeachFunc(model *C.GtkTreeModel, path *C.GtkTreePath, iter *C.GtkTreeIter, data C.gpointer) C.gboolean {
	id := int(uintptr(data))

	treeModelForeachFuncRegistry.Lock()
	r := treeModelForeachFuncRegistry.m[id]
	treeModelForeachFuncRegistry.Unlock()

	goPath := &TreePath{(*C.GtkTreePath)(path)}
	goIter := &TreeIter{(C.GtkTreeIter)(*iter)}
	return gbool(r.fn(
		wrapTreeModel(glib.Take(unsafe.Pointer(model))),
		goPath,
		goIter,
		r.userData))
}
