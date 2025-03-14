// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Jobs() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"bg-white text-black p-4 mb-4\"><h1 class=\"text-3xl mb-4\">Jobs</h1><h2 class=\"text-2xl mb-4\">Start a job:</h2><ul><li class=\"mb-4\"><h3 class=\"text-lg mb-4\">Screenshot</h3><button id=\"startJobButton\" class=\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded\" hx-post=\"/admin/jobs/start/screenshot-all\" hx-confirm=\"Are you sure you want to generate screenshots for all documents?\">Screenshot All documents</button></li><li class=\"mb-4\"><h3 class=\"text-lg mb-4\">Snapshot</h3><form hx-post=\"/admin/jobs/start/snapshot\" hx-target=\"#snapshot-response-target\"><input type=\"text\" name=\"document_id\" placeholder=\"Document ID\"> <button id=\"startSnapshotButton\" class=\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded\" type=\"submit\">Start snapshot</button></form><div id=\"snapshot-response-target\"></div><form hx-post=\"/admin/jobs/start/snapshot-all\" hx-target=\"#snapshot-all-response-target\"><input type=\"text\" name=\"version\" placeholder=\"Version\"> <button id=\"startSnapshotAllButton\" class=\"bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded\" type=\"submit\">Start snapshot all</button></form><div id=\"snapshot-all-response-target\"></div></li></ul></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = AdminLayout("jobs").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
