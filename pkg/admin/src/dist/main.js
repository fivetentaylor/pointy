"use strict";
(() => {
  var __create = Object.create;
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getProtoOf = Object.getPrototypeOf, __hasOwnProp = Object.prototype.hasOwnProperty;
  var __commonJS = (cb, mod) => function() {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from == "object" || typeof from == "function")
      for (let key of __getOwnPropNames(from))
        !__hasOwnProp.call(to, key) && key !== except && __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    return to;
  };
  var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
    // If the importer is in node compatibility mode or this is not an ESM
    // file that has been converted to a CommonJS file using a Babel-
    // compatible transform (i.e. "__esModule" has not been set), then set
    // "default" to the CommonJS "module.exports" for node compatibility.
    isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: !0 }) : target,
    mod
  ));

  // pkg/admin/src/node_modules/htmx.org/dist/htmx.min.js
  var require_htmx_min = __commonJS({
    "pkg/admin/src/node_modules/htmx.org/dist/htmx.min.js"(exports, module) {
      (function(e2, t2) {
        typeof define == "function" && define.amd ? define([], t2) : typeof module == "object" && module.exports ? module.exports = t2() : e2.htmx = e2.htmx || t2();
      })(typeof self < "u" ? self : exports, function() {
        return function() {
          "use strict";
          var Q = { onLoad: F, process: zt, on: de, off: ge, trigger: ce, ajax: Nr, find: C, findAll: f, closest: v, values: function(e2, t2) {
            var r2 = dr(e2, t2 || "post");
            return r2.values;
          }, remove: _, addClass: z, removeClass: n, toggleClass: $, takeClass: W, defineExtension: Ur, removeExtension: Br, logAll: V, logNone: j, logger: null, config: { historyEnabled: !0, historyCacheSize: 10, refreshOnHistoryMiss: !1, defaultSwapStyle: "innerHTML", defaultSwapDelay: 0, defaultSettleDelay: 20, includeIndicatorStyles: !0, indicatorClass: "htmx-indicator", requestClass: "htmx-request", addedClass: "htmx-added", settlingClass: "htmx-settling", swappingClass: "htmx-swapping", allowEval: !0, allowScriptTags: !0, inlineScriptNonce: "", attributesToSettle: ["class", "style", "width", "height"], withCredentials: !1, timeout: 0, wsReconnectDelay: "full-jitter", wsBinaryType: "blob", disableSelector: "[hx-disable], [data-hx-disable]", useTemplateFragments: !1, scrollBehavior: "smooth", defaultFocusScroll: !1, getCacheBusterParam: !1, globalViewTransitions: !1, methodsThatUseUrlParams: ["get"], selfRequestsOnly: !1, ignoreTitle: !1, scrollIntoViewOnBoost: !0, triggerSpecsCache: null }, parseInterval: d, _: t, createEventSource: function(e2) {
            return new EventSource(e2, { withCredentials: !0 });
          }, createWebSocket: function(e2) {
            var t2 = new WebSocket(e2, []);
            return t2.binaryType = Q.config.wsBinaryType, t2;
          }, version: "1.9.10" }, r = { addTriggerHandler: Lt, bodyContains: se, canAccessLocalStorage: U, findThisElement: xe, filterValues: yr, hasAttribute: o, getAttributeValue: te, getClosestAttributeValue: ne, getClosestMatch: c, getExpressionVars: Hr, getHeaders: xr, getInputValues: dr, getInternalData: ae, getSwapSpecification: wr, getTriggerSpecs: it, getTarget: ye, makeFragment: l, mergeObjects: le, makeSettleInfo: T, oobSwap: Ee, querySelectorExt: ue, selectAndSwap: je, settleImmediately: nr, shouldCancel: ut, triggerEvent: ce, triggerErrorEvent: fe, withExtensions: R }, w = ["get", "post", "put", "delete", "patch"], i = w.map(function(e2) {
            return "[hx-" + e2 + "], [data-hx-" + e2 + "]";
          }).join(", "), S = e("head"), q = e("title"), H = e("svg", !0);
          function e(e2, t2 = !1) {
            return new RegExp(`<${e2}(\\s[^>]*>|>)([\\s\\S]*?)<\\/${e2}>`, t2 ? "gim" : "im");
          }
          function d(e2) {
            if (e2 == null)
              return;
            let t2 = NaN;
            return e2.slice(-2) == "ms" ? t2 = parseFloat(e2.slice(0, -2)) : e2.slice(-1) == "s" ? t2 = parseFloat(e2.slice(0, -1)) * 1e3 : e2.slice(-1) == "m" ? t2 = parseFloat(e2.slice(0, -1)) * 1e3 * 60 : t2 = parseFloat(e2), isNaN(t2) ? void 0 : t2;
          }
          function ee(e2, t2) {
            return e2.getAttribute && e2.getAttribute(t2);
          }
          function o(e2, t2) {
            return e2.hasAttribute && (e2.hasAttribute(t2) || e2.hasAttribute("data-" + t2));
          }
          function te(e2, t2) {
            return ee(e2, t2) || ee(e2, "data-" + t2);
          }
          function u(e2) {
            return e2.parentElement;
          }
          function re() {
            return document;
          }
          function c(e2, t2) {
            for (; e2 && !t2(e2); )
              e2 = u(e2);
            return e2 || null;
          }
          function L(e2, t2, r2) {
            var n2 = te(t2, r2), i2 = te(t2, "hx-disinherit");
            return e2 !== t2 && i2 && (i2 === "*" || i2.split(" ").indexOf(r2) >= 0) ? "unset" : n2;
          }
          function ne(t2, r2) {
            var n2 = null;
            if (c(t2, function(e2) {
              return n2 = L(t2, e2, r2);
            }), n2 !== "unset")
              return n2;
          }
          function h(e2, t2) {
            var r2 = e2.matches || e2.matchesSelector || e2.msMatchesSelector || e2.mozMatchesSelector || e2.webkitMatchesSelector || e2.oMatchesSelector;
            return r2 && r2.call(e2, t2);
          }
          function A(e2) {
            var t2 = /<([a-z][^\/\0>\x20\t\r\n\f]*)/i, r2 = t2.exec(e2);
            return r2 ? r2[1].toLowerCase() : "";
          }
          function a(e2, t2) {
            for (var r2 = new DOMParser(), n2 = r2.parseFromString(e2, "text/html"), i2 = n2.body; t2 > 0; )
              t2--, i2 = i2.firstChild;
            return i2 == null && (i2 = re().createDocumentFragment()), i2;
          }
          function N(e2) {
            return /<body/.test(e2);
          }
          function l(e2) {
            var t2 = !N(e2), r2 = A(e2), n2 = e2;
            if (r2 === "head" && (n2 = n2.replace(S, "")), Q.config.useTemplateFragments && t2) {
              var i2 = a("<body><template>" + n2 + "</template></body>", 0);
              return i2.querySelector("template").content;
            }
            switch (r2) {
              case "thead":
              case "tbody":
              case "tfoot":
              case "colgroup":
              case "caption":
                return a("<table>" + n2 + "</table>", 1);
              case "col":
                return a("<table><colgroup>" + n2 + "</colgroup></table>", 2);
              case "tr":
                return a("<table><tbody>" + n2 + "</tbody></table>", 2);
              case "td":
              case "th":
                return a("<table><tbody><tr>" + n2 + "</tr></tbody></table>", 3);
              case "script":
              case "style":
                return a("<div>" + n2 + "</div>", 1);
              default:
                return a(n2, 0);
            }
          }
          function ie(e2) {
            e2 && e2();
          }
          function I(e2, t2) {
            return Object.prototype.toString.call(e2) === "[object " + t2 + "]";
          }
          function k(e2) {
            return I(e2, "Function");
          }
          function P(e2) {
            return I(e2, "Object");
          }
          function ae(e2) {
            var t2 = "htmx-internal-data", r2 = e2[t2];
            return r2 || (r2 = e2[t2] = {}), r2;
          }
          function M(e2) {
            var t2 = [];
            if (e2)
              for (var r2 = 0; r2 < e2.length; r2++)
                t2.push(e2[r2]);
            return t2;
          }
          function oe(e2, t2) {
            if (e2)
              for (var r2 = 0; r2 < e2.length; r2++)
                t2(e2[r2]);
          }
          function X(e2) {
            var t2 = e2.getBoundingClientRect(), r2 = t2.top, n2 = t2.bottom;
            return r2 < window.innerHeight && n2 >= 0;
          }
          function se(e2) {
            return e2.getRootNode && e2.getRootNode() instanceof window.ShadowRoot ? re().body.contains(e2.getRootNode().host) : re().body.contains(e2);
          }
          function D(e2) {
            return e2.trim().split(/\s+/);
          }
          function le(e2, t2) {
            for (var r2 in t2)
              t2.hasOwnProperty(r2) && (e2[r2] = t2[r2]);
            return e2;
          }
          function E(e2) {
            try {
              return JSON.parse(e2);
            } catch (e3) {
              return b(e3), null;
            }
          }
          function U() {
            var e2 = "htmx:localStorageTest";
            try {
              return localStorage.setItem(e2, e2), localStorage.removeItem(e2), !0;
            } catch {
              return !1;
            }
          }
          function B(t2) {
            try {
              var e2 = new URL(t2);
              return e2 && (t2 = e2.pathname + e2.search), /^\/$/.test(t2) || (t2 = t2.replace(/\/+$/, "")), t2;
            } catch {
              return t2;
            }
          }
          function t(e) {
            return Tr(re().body, function() {
              return eval(e);
            });
          }
          function F(t2) {
            var e2 = Q.on("htmx:load", function(e3) {
              t2(e3.detail.elt);
            });
            return e2;
          }
          function V() {
            Q.logger = function(e2, t2, r2) {
              console && console.log(t2, e2, r2);
            };
          }
          function j() {
            Q.logger = null;
          }
          function C(e2, t2) {
            return t2 ? e2.querySelector(t2) : C(re(), e2);
          }
          function f(e2, t2) {
            return t2 ? e2.querySelectorAll(t2) : f(re(), e2);
          }
          function _(e2, t2) {
            e2 = g(e2), t2 ? setTimeout(function() {
              _(e2), e2 = null;
            }, t2) : e2.parentElement.removeChild(e2);
          }
          function z(e2, t2, r2) {
            e2 = g(e2), r2 ? setTimeout(function() {
              z(e2, t2), e2 = null;
            }, r2) : e2.classList && e2.classList.add(t2);
          }
          function n(e2, t2, r2) {
            e2 = g(e2), r2 ? setTimeout(function() {
              n(e2, t2), e2 = null;
            }, r2) : e2.classList && (e2.classList.remove(t2), e2.classList.length === 0 && e2.removeAttribute("class"));
          }
          function $(e2, t2) {
            e2 = g(e2), e2.classList.toggle(t2);
          }
          function W(e2, t2) {
            e2 = g(e2), oe(e2.parentElement.children, function(e3) {
              n(e3, t2);
            }), z(e2, t2);
          }
          function v(e2, t2) {
            if (e2 = g(e2), e2.closest)
              return e2.closest(t2);
            do
              if (e2 == null || h(e2, t2))
                return e2;
            while (e2 = e2 && u(e2));
            return null;
          }
          function s(e2, t2) {
            return e2.substring(0, t2.length) === t2;
          }
          function G(e2, t2) {
            return e2.substring(e2.length - t2.length) === t2;
          }
          function J(e2) {
            var t2 = e2.trim();
            return s(t2, "<") && G(t2, "/>") ? t2.substring(1, t2.length - 2) : t2;
          }
          function Z(e2, t2) {
            return t2.indexOf("closest ") === 0 ? [v(e2, J(t2.substr(8)))] : t2.indexOf("find ") === 0 ? [C(e2, J(t2.substr(5)))] : t2 === "next" ? [e2.nextElementSibling] : t2.indexOf("next ") === 0 ? [K(e2, J(t2.substr(5)))] : t2 === "previous" ? [e2.previousElementSibling] : t2.indexOf("previous ") === 0 ? [Y(e2, J(t2.substr(9)))] : t2 === "document" ? [document] : t2 === "window" ? [window] : t2 === "body" ? [document.body] : re().querySelectorAll(J(t2));
          }
          var K = function(e2, t2) {
            for (var r2 = re().querySelectorAll(t2), n2 = 0; n2 < r2.length; n2++) {
              var i2 = r2[n2];
              if (i2.compareDocumentPosition(e2) === Node.DOCUMENT_POSITION_PRECEDING)
                return i2;
            }
          }, Y = function(e2, t2) {
            for (var r2 = re().querySelectorAll(t2), n2 = r2.length - 1; n2 >= 0; n2--) {
              var i2 = r2[n2];
              if (i2.compareDocumentPosition(e2) === Node.DOCUMENT_POSITION_FOLLOWING)
                return i2;
            }
          };
          function ue(e2, t2) {
            return t2 ? Z(e2, t2)[0] : Z(re().body, e2)[0];
          }
          function g(e2) {
            return I(e2, "String") ? C(e2) : e2;
          }
          function ve(e2, t2, r2) {
            return k(t2) ? { target: re().body, event: e2, listener: t2 } : { target: g(e2), event: t2, listener: r2 };
          }
          function de(t2, r2, n2) {
            jr(function() {
              var e3 = ve(t2, r2, n2);
              e3.target.addEventListener(e3.event, e3.listener);
            });
            var e2 = k(r2);
            return e2 ? r2 : n2;
          }
          function ge(t2, r2, n2) {
            return jr(function() {
              var e2 = ve(t2, r2, n2);
              e2.target.removeEventListener(e2.event, e2.listener);
            }), k(r2) ? r2 : n2;
          }
          var me = re().createElement("output");
          function pe(e2, t2) {
            var r2 = ne(e2, t2);
            if (r2) {
              if (r2 === "this")
                return [xe(e2, t2)];
              var n2 = Z(e2, r2);
              return n2.length === 0 ? (b('The selector "' + r2 + '" on ' + t2 + " returned no matches!"), [me]) : n2;
            }
          }
          function xe(e2, t2) {
            return c(e2, function(e3) {
              return te(e3, t2) != null;
            });
          }
          function ye(e2) {
            var t2 = ne(e2, "hx-target");
            if (t2)
              return t2 === "this" ? xe(e2, "hx-target") : ue(e2, t2);
            var r2 = ae(e2);
            return r2.boosted ? re().body : e2;
          }
          function be(e2) {
            for (var t2 = Q.config.attributesToSettle, r2 = 0; r2 < t2.length; r2++)
              if (e2 === t2[r2])
                return !0;
            return !1;
          }
          function we(t2, r2) {
            oe(t2.attributes, function(e2) {
              !r2.hasAttribute(e2.name) && be(e2.name) && t2.removeAttribute(e2.name);
            }), oe(r2.attributes, function(e2) {
              be(e2.name) && t2.setAttribute(e2.name, e2.value);
            });
          }
          function Se(e2, t2) {
            for (var r2 = Fr(t2), n2 = 0; n2 < r2.length; n2++) {
              var i2 = r2[n2];
              try {
                if (i2.isInlineSwap(e2))
                  return !0;
              } catch (e3) {
                b(e3);
              }
            }
            return e2 === "outerHTML";
          }
          function Ee(e2, i2, a2) {
            var t2 = "#" + ee(i2, "id"), o2 = "outerHTML";
            e2 === "true" || (e2.indexOf(":") > 0 ? (o2 = e2.substr(0, e2.indexOf(":")), t2 = e2.substr(e2.indexOf(":") + 1, e2.length)) : o2 = e2);
            var r2 = re().querySelectorAll(t2);
            return r2 ? (oe(r2, function(e3) {
              var t3, r3 = i2.cloneNode(!0);
              t3 = re().createDocumentFragment(), t3.appendChild(r3), Se(o2, e3) || (t3 = r3);
              var n2 = { shouldSwap: !0, target: e3, fragment: t3 };
              ce(e3, "htmx:oobBeforeSwap", n2) && (e3 = n2.target, n2.shouldSwap && Fe(o2, e3, e3, t3, a2), oe(a2.elts, function(e4) {
                ce(e4, "htmx:oobAfterSwap", n2);
              }));
            }), i2.parentNode.removeChild(i2)) : (i2.parentNode.removeChild(i2), fe(re().body, "htmx:oobErrorNoTarget", { content: i2 })), e2;
          }
          function Ce(e2, t2, r2) {
            var n2 = ne(e2, "hx-select-oob");
            if (n2)
              for (var i2 = n2.split(","), a2 = 0; a2 < i2.length; a2++) {
                var o2 = i2[a2].split(":", 2), s2 = o2[0].trim();
                s2.indexOf("#") === 0 && (s2 = s2.substring(1));
                var l2 = o2[1] || "true", u2 = t2.querySelector("#" + s2);
                u2 && Ee(l2, u2, r2);
              }
            oe(f(t2, "[hx-swap-oob], [data-hx-swap-oob]"), function(e3) {
              var t3 = te(e3, "hx-swap-oob");
              t3 != null && Ee(t3, e3, r2);
            });
          }
          function Re(e2) {
            oe(f(e2, "[hx-preserve], [data-hx-preserve]"), function(e3) {
              var t2 = te(e3, "id"), r2 = re().getElementById(t2);
              r2 != null && e3.parentNode.replaceChild(r2, e3);
            });
          }
          function Te(o2, e2, s2) {
            oe(e2.querySelectorAll("[id]"), function(e3) {
              var t2 = ee(e3, "id");
              if (t2 && t2.length > 0) {
                var r2 = t2.replace("'", "\\'"), n2 = e3.tagName.replace(":", "\\:"), i2 = o2.querySelector(n2 + "[id='" + r2 + "']");
                if (i2 && i2 !== o2) {
                  var a2 = e3.cloneNode();
                  we(e3, i2), s2.tasks.push(function() {
                    we(e3, a2);
                  });
                }
              }
            });
          }
          function Oe(e2) {
            return function() {
              n(e2, Q.config.addedClass), zt(e2), Nt(e2), qe(e2), ce(e2, "htmx:load");
            };
          }
          function qe(e2) {
            var t2 = "[autofocus]", r2 = h(e2, t2) ? e2 : e2.querySelector(t2);
            r2?.focus();
          }
          function m(e2, t2, r2, n2) {
            for (Te(e2, r2, n2); r2.childNodes.length > 0; ) {
              var i2 = r2.firstChild;
              z(i2, Q.config.addedClass), e2.insertBefore(i2, t2), i2.nodeType !== Node.TEXT_NODE && i2.nodeType !== Node.COMMENT_NODE && n2.tasks.push(Oe(i2));
            }
          }
          function He(e2, t2) {
            for (var r2 = 0; r2 < e2.length; )
              t2 = (t2 << 5) - t2 + e2.charCodeAt(r2++) | 0;
            return t2;
          }
          function Le(e2) {
            var t2 = 0;
            if (e2.attributes)
              for (var r2 = 0; r2 < e2.attributes.length; r2++) {
                var n2 = e2.attributes[r2];
                n2.value && (t2 = He(n2.name, t2), t2 = He(n2.value, t2));
              }
            return t2;
          }
          function Ae(e2) {
            var t2 = ae(e2);
            if (t2.onHandlers) {
              for (var r2 = 0; r2 < t2.onHandlers.length; r2++) {
                let n2 = t2.onHandlers[r2];
                e2.removeEventListener(n2.event, n2.listener);
              }
              delete t2.onHandlers;
            }
          }
          function Ne(e2) {
            var t2 = ae(e2);
            t2.timeout && clearTimeout(t2.timeout), t2.webSocket && t2.webSocket.close(), t2.sseEventSource && t2.sseEventSource.close(), t2.listenerInfos && oe(t2.listenerInfos, function(e3) {
              e3.on && e3.on.removeEventListener(e3.trigger, e3.listener);
            }), Ae(e2), oe(Object.keys(t2), function(e3) {
              delete t2[e3];
            });
          }
          function p(e2) {
            ce(e2, "htmx:beforeCleanupElement"), Ne(e2), e2.children && oe(e2.children, function(e3) {
              p(e3);
            });
          }
          function Ie(t2, e2, r2) {
            if (t2.tagName === "BODY")
              return Ue(t2, e2, r2);
            var n2, i2 = t2.previousSibling;
            for (m(u(t2), t2, e2, r2), i2 == null ? n2 = u(t2).firstChild : n2 = i2.nextSibling, r2.elts = r2.elts.filter(function(e3) {
              return e3 != t2;
            }); n2 && n2 !== t2; )
              n2.nodeType === Node.ELEMENT_NODE && r2.elts.push(n2), n2 = n2.nextElementSibling;
            p(t2), u(t2).removeChild(t2);
          }
          function ke(e2, t2, r2) {
            return m(e2, e2.firstChild, t2, r2);
          }
          function Pe(e2, t2, r2) {
            return m(u(e2), e2, t2, r2);
          }
          function Me(e2, t2, r2) {
            return m(e2, null, t2, r2);
          }
          function Xe(e2, t2, r2) {
            return m(u(e2), e2.nextSibling, t2, r2);
          }
          function De(e2, t2, r2) {
            return p(e2), u(e2).removeChild(e2);
          }
          function Ue(e2, t2, r2) {
            var n2 = e2.firstChild;
            if (m(e2, n2, t2, r2), n2) {
              for (; n2.nextSibling; )
                p(n2.nextSibling), e2.removeChild(n2.nextSibling);
              p(n2), e2.removeChild(n2);
            }
          }
          function Be(e2, t2, r2) {
            var n2 = r2 || ne(e2, "hx-select");
            if (n2) {
              var i2 = re().createDocumentFragment();
              oe(t2.querySelectorAll(n2), function(e3) {
                i2.appendChild(e3);
              }), t2 = i2;
            }
            return t2;
          }
          function Fe(e2, t2, r2, n2, i2) {
            switch (e2) {
              case "none":
                return;
              case "outerHTML":
                Ie(r2, n2, i2);
                return;
              case "afterbegin":
                ke(r2, n2, i2);
                return;
              case "beforebegin":
                Pe(r2, n2, i2);
                return;
              case "beforeend":
                Me(r2, n2, i2);
                return;
              case "afterend":
                Xe(r2, n2, i2);
                return;
              case "delete":
                De(r2, n2, i2);
                return;
              default:
                for (var a2 = Fr(t2), o2 = 0; o2 < a2.length; o2++) {
                  var s2 = a2[o2];
                  try {
                    var l2 = s2.handleSwap(e2, r2, n2, i2);
                    if (l2) {
                      if (typeof l2.length < "u")
                        for (var u2 = 0; u2 < l2.length; u2++) {
                          var f2 = l2[u2];
                          f2.nodeType !== Node.TEXT_NODE && f2.nodeType !== Node.COMMENT_NODE && i2.tasks.push(Oe(f2));
                        }
                      return;
                    }
                  } catch (e3) {
                    b(e3);
                  }
                }
                e2 === "innerHTML" ? Ue(r2, n2, i2) : Fe(Q.config.defaultSwapStyle, t2, r2, n2, i2);
            }
          }
          function Ve(e2) {
            if (e2.indexOf("<title") > -1) {
              var t2 = e2.replace(H, ""), r2 = t2.match(q);
              if (r2)
                return r2[2];
            }
          }
          function je(e2, t2, r2, n2, i2, a2) {
            i2.title = Ve(n2);
            var o2 = l(n2);
            if (o2)
              return Ce(r2, o2, i2), o2 = Be(r2, o2, a2), Re(o2), Fe(e2, r2, t2, o2, i2);
          }
          function _e(e2, t2, r2) {
            var n2 = e2.getResponseHeader(t2);
            if (n2.indexOf("{") === 0) {
              var i2 = E(n2);
              for (var a2 in i2)
                if (i2.hasOwnProperty(a2)) {
                  var o2 = i2[a2];
                  P(o2) || (o2 = { value: o2 }), ce(r2, a2, o2);
                }
            } else
              for (var s2 = n2.split(","), l2 = 0; l2 < s2.length; l2++)
                ce(r2, s2[l2].trim(), []);
          }
          var ze = /\s/, x = /[\s,]/, $e = /[_$a-zA-Z]/, We = /[_$a-zA-Z0-9]/, Ge = ['"', "'", "/"], Je = /[^\s]/, Ze = /[{(]/, Ke = /[})]/;
          function Ye(e2) {
            for (var t2 = [], r2 = 0; r2 < e2.length; ) {
              if ($e.exec(e2.charAt(r2))) {
                for (var n2 = r2; We.exec(e2.charAt(r2 + 1)); )
                  r2++;
                t2.push(e2.substr(n2, r2 - n2 + 1));
              } else if (Ge.indexOf(e2.charAt(r2)) !== -1) {
                var i2 = e2.charAt(r2), n2 = r2;
                for (r2++; r2 < e2.length && e2.charAt(r2) !== i2; )
                  e2.charAt(r2) === "\\" && r2++, r2++;
                t2.push(e2.substr(n2, r2 - n2 + 1));
              } else {
                var a2 = e2.charAt(r2);
                t2.push(a2);
              }
              r2++;
            }
            return t2;
          }
          function Qe(e2, t2, r2) {
            return $e.exec(e2.charAt(0)) && e2 !== "true" && e2 !== "false" && e2 !== "this" && e2 !== r2 && t2 !== ".";
          }
          function et(e2, t2, r2) {
            if (t2[0] === "[") {
              t2.shift();
              for (var n2 = 1, i2 = " return (function(" + r2 + "){ return (", a2 = null; t2.length > 0; ) {
                var o2 = t2[0];
                if (o2 === "]") {
                  if (n2--, n2 === 0) {
                    a2 === null && (i2 = i2 + "true"), t2.shift(), i2 += ")})";
                    try {
                      var s2 = Tr(e2, function() {
                        return Function(i2)();
                      }, function() {
                        return !0;
                      });
                      return s2.source = i2, s2;
                    } catch (e3) {
                      return fe(re().body, "htmx:syntax:error", { error: e3, source: i2 }), null;
                    }
                  }
                } else o2 === "[" && n2++;
                Qe(o2, a2, r2) ? i2 += "((" + r2 + "." + o2 + ") ? (" + r2 + "." + o2 + ") : (window." + o2 + "))" : i2 = i2 + o2, a2 = t2.shift();
              }
            }
          }
          function y(e2, t2) {
            for (var r2 = ""; e2.length > 0 && !t2.test(e2[0]); )
              r2 += e2.shift();
            return r2;
          }
          function tt(e2) {
            var t2;
            return e2.length > 0 && Ze.test(e2[0]) ? (e2.shift(), t2 = y(e2, Ke).trim(), e2.shift()) : t2 = y(e2, x), t2;
          }
          var rt = "input, textarea, select";
          function nt(e2, t2, r2) {
            var n2 = [], i2 = Ye(t2);
            do {
              y(i2, Je);
              var a2 = i2.length, o2 = y(i2, /[,\[\s]/);
              if (o2 !== "")
                if (o2 === "every") {
                  var s2 = { trigger: "every" };
                  y(i2, Je), s2.pollInterval = d(y(i2, /[,\[\s]/)), y(i2, Je);
                  var l2 = et(e2, i2, "event");
                  l2 && (s2.eventFilter = l2), n2.push(s2);
                } else if (o2.indexOf("sse:") === 0)
                  n2.push({ trigger: "sse", sseEvent: o2.substr(4) });
                else {
                  var u2 = { trigger: o2 }, l2 = et(e2, i2, "event");
                  for (l2 && (u2.eventFilter = l2); i2.length > 0 && i2[0] !== ","; ) {
                    y(i2, Je);
                    var f2 = i2.shift();
                    if (f2 === "changed")
                      u2.changed = !0;
                    else if (f2 === "once")
                      u2.once = !0;
                    else if (f2 === "consume")
                      u2.consume = !0;
                    else if (f2 === "delay" && i2[0] === ":")
                      i2.shift(), u2.delay = d(y(i2, x));
                    else if (f2 === "from" && i2[0] === ":") {
                      if (i2.shift(), Ze.test(i2[0]))
                        var c2 = tt(i2);
                      else {
                        var c2 = y(i2, x);
                        if (c2 === "closest" || c2 === "find" || c2 === "next" || c2 === "previous") {
                          i2.shift();
                          var h2 = tt(i2);
                          h2.length > 0 && (c2 += " " + h2);
                        }
                      }
                      u2.from = c2;
                    } else f2 === "target" && i2[0] === ":" ? (i2.shift(), u2.target = tt(i2)) : f2 === "throttle" && i2[0] === ":" ? (i2.shift(), u2.throttle = d(y(i2, x))) : f2 === "queue" && i2[0] === ":" ? (i2.shift(), u2.queue = y(i2, x)) : f2 === "root" && i2[0] === ":" ? (i2.shift(), u2[f2] = tt(i2)) : f2 === "threshold" && i2[0] === ":" ? (i2.shift(), u2[f2] = y(i2, x)) : fe(e2, "htmx:syntax:error", { token: i2.shift() });
                  }
                  n2.push(u2);
                }
              i2.length === a2 && fe(e2, "htmx:syntax:error", { token: i2.shift() }), y(i2, Je);
            } while (i2[0] === "," && i2.shift());
            return r2 && (r2[t2] = n2), n2;
          }
          function it(e2) {
            var t2 = te(e2, "hx-trigger"), r2 = [];
            if (t2) {
              var n2 = Q.config.triggerSpecsCache;
              r2 = n2 && n2[t2] || nt(e2, t2, n2);
            }
            return r2.length > 0 ? r2 : h(e2, "form") ? [{ trigger: "submit" }] : h(e2, 'input[type="button"], input[type="submit"]') ? [{ trigger: "click" }] : h(e2, rt) ? [{ trigger: "change" }] : [{ trigger: "click" }];
          }
          function at(e2) {
            ae(e2).cancelled = !0;
          }
          function ot(e2, t2, r2) {
            var n2 = ae(e2);
            n2.timeout = setTimeout(function() {
              se(e2) && n2.cancelled !== !0 && (ct(r2, e2, Wt("hx:poll:trigger", { triggerSpec: r2, target: e2 })) || t2(e2), ot(e2, t2, r2));
            }, r2.pollInterval);
          }
          function st(e2) {
            return location.hostname === e2.hostname && ee(e2, "href") && ee(e2, "href").indexOf("#") !== 0;
          }
          function lt(t2, r2, e2) {
            if (t2.tagName === "A" && st(t2) && (t2.target === "" || t2.target === "_self") || t2.tagName === "FORM") {
              r2.boosted = !0;
              var n2, i2;
              if (t2.tagName === "A")
                n2 = "get", i2 = ee(t2, "href");
              else {
                var a2 = ee(t2, "method");
                n2 = a2 ? a2.toLowerCase() : "get", i2 = ee(t2, "action");
              }
              e2.forEach(function(e3) {
                ht(t2, function(e4, t3) {
                  if (v(e4, Q.config.disableSelector)) {
                    p(e4);
                    return;
                  }
                  he(n2, i2, e4, t3);
                }, r2, e3, !0);
              });
            }
          }
          function ut(e2, t2) {
            return !!((e2.type === "submit" || e2.type === "click") && (t2.tagName === "FORM" || h(t2, 'input[type="submit"], button') && v(t2, "form") !== null || t2.tagName === "A" && t2.href && (t2.getAttribute("href") === "#" || t2.getAttribute("href").indexOf("#") !== 0)));
          }
          function ft(e2, t2) {
            return ae(e2).boosted && e2.tagName === "A" && t2.type === "click" && (t2.ctrlKey || t2.metaKey);
          }
          function ct(e2, t2, r2) {
            var n2 = e2.eventFilter;
            if (n2)
              try {
                return n2.call(t2, r2) !== !0;
              } catch (e3) {
                return fe(re().body, "htmx:eventFilter:error", { error: e3, source: n2.source }), !0;
              }
            return !1;
          }
          function ht(a2, o2, e2, s2, l2) {
            var u2 = ae(a2), t2;
            s2.from ? t2 = Z(a2, s2.from) : t2 = [a2], s2.changed && t2.forEach(function(e3) {
              var t3 = ae(e3);
              t3.lastValue = e3.value;
            }), oe(t2, function(n2) {
              var i2 = function(e3) {
                if (!se(a2)) {
                  n2.removeEventListener(s2.trigger, i2);
                  return;
                }
                if (!ft(a2, e3) && ((l2 || ut(e3, a2)) && e3.preventDefault(), !ct(s2, a2, e3))) {
                  var t3 = ae(e3);
                  if (t3.triggerSpec = s2, t3.handledFor == null && (t3.handledFor = []), t3.handledFor.indexOf(a2) < 0) {
                    if (t3.handledFor.push(a2), s2.consume && e3.stopPropagation(), s2.target && e3.target && !h(e3.target, s2.target))
                      return;
                    if (s2.once) {
                      if (u2.triggeredOnce)
                        return;
                      u2.triggeredOnce = !0;
                    }
                    if (s2.changed) {
                      var r2 = ae(n2);
                      if (r2.lastValue === n2.value)
                        return;
                      r2.lastValue = n2.value;
                    }
                    if (u2.delayed && clearTimeout(u2.delayed), u2.throttle)
                      return;
                    s2.throttle > 0 ? u2.throttle || (o2(a2, e3), u2.throttle = setTimeout(function() {
                      u2.throttle = null;
                    }, s2.throttle)) : s2.delay > 0 ? u2.delayed = setTimeout(function() {
                      o2(a2, e3);
                    }, s2.delay) : (ce(a2, "htmx:trigger"), o2(a2, e3));
                  }
                }
              };
              e2.listenerInfos == null && (e2.listenerInfos = []), e2.listenerInfos.push({ trigger: s2.trigger, listener: i2, on: n2 }), n2.addEventListener(s2.trigger, i2);
            });
          }
          var vt = !1, dt = null;
          function gt() {
            dt || (dt = function() {
              vt = !0;
            }, window.addEventListener("scroll", dt), setInterval(function() {
              vt && (vt = !1, oe(re().querySelectorAll("[hx-trigger='revealed'],[data-hx-trigger='revealed']"), function(e2) {
                mt(e2);
              }));
            }, 200));
          }
          function mt(t2) {
            if (!o(t2, "data-hx-revealed") && X(t2)) {
              t2.setAttribute("data-hx-revealed", "true");
              var e2 = ae(t2);
              e2.initHash ? ce(t2, "revealed") : t2.addEventListener("htmx:afterProcessNode", function(e3) {
                ce(t2, "revealed");
              }, { once: !0 });
            }
          }
          function pt(e2, t2, r2) {
            for (var n2 = D(r2), i2 = 0; i2 < n2.length; i2++) {
              var a2 = n2[i2].split(/:(.+)/);
              a2[0] === "connect" && xt(e2, a2[1], 0), a2[0] === "send" && bt(e2);
            }
          }
          function xt(s2, r2, n2) {
            if (se(s2)) {
              if (r2.indexOf("/") == 0) {
                var e2 = location.hostname + (location.port ? ":" + location.port : "");
                location.protocol == "https:" ? r2 = "wss://" + e2 + r2 : location.protocol == "http:" && (r2 = "ws://" + e2 + r2);
              }
              var t2 = Q.createWebSocket(r2);
              t2.onerror = function(e3) {
                fe(s2, "htmx:wsError", { error: e3, socket: t2 }), yt(s2);
              }, t2.onclose = function(e3) {
                if ([1006, 1012, 1013].indexOf(e3.code) >= 0) {
                  var t3 = wt(n2);
                  setTimeout(function() {
                    xt(s2, r2, n2 + 1);
                  }, t3);
                }
              }, t2.onopen = function(e3) {
                n2 = 0;
              }, ae(s2).webSocket = t2, t2.addEventListener("message", function(e3) {
                if (!yt(s2)) {
                  var t3 = e3.data;
                  R(s2, function(e4) {
                    t3 = e4.transformResponse(t3, null, s2);
                  });
                  for (var r3 = T(s2), n3 = l(t3), i2 = M(n3.children), a2 = 0; a2 < i2.length; a2++) {
                    var o2 = i2[a2];
                    Ee(te(o2, "hx-swap-oob") || "true", o2, r3);
                  }
                  nr(r3.tasks);
                }
              });
            }
          }
          function yt(e2) {
            if (!se(e2))
              return ae(e2).webSocket.close(), !0;
          }
          function bt(u2) {
            var f2 = c(u2, function(e2) {
              return ae(e2).webSocket != null;
            });
            f2 ? u2.addEventListener(it(u2)[0].trigger, function(e2) {
              var t2 = ae(f2).webSocket, r2 = xr(u2, f2), n2 = dr(u2, "post"), i2 = n2.errors, a2 = n2.values, o2 = Hr(u2), s2 = le(a2, o2), l2 = yr(s2, u2);
              if (l2.HEADERS = r2, i2 && i2.length > 0) {
                ce(u2, "htmx:validation:halted", i2);
                return;
              }
              t2.send(JSON.stringify(l2)), ut(e2, u2) && e2.preventDefault();
            }) : fe(u2, "htmx:noWebSocketSourceError");
          }
          function wt(e2) {
            var t2 = Q.config.wsReconnectDelay;
            if (typeof t2 == "function")
              return t2(e2);
            if (t2 === "full-jitter") {
              var r2 = Math.min(e2, 6), n2 = 1e3 * Math.pow(2, r2);
              return n2 * Math.random();
            }
            b('htmx.config.wsReconnectDelay must either be a function or the string "full-jitter"');
          }
          function St(e2, t2, r2) {
            for (var n2 = D(r2), i2 = 0; i2 < n2.length; i2++) {
              var a2 = n2[i2].split(/:(.+)/);
              a2[0] === "connect" && Et(e2, a2[1]), a2[0] === "swap" && Ct(e2, a2[1]);
            }
          }
          function Et(t2, e2) {
            var r2 = Q.createEventSource(e2);
            r2.onerror = function(e3) {
              fe(t2, "htmx:sseError", { error: e3, source: r2 }), Tt(t2);
            }, ae(t2).sseEventSource = r2;
          }
          function Ct(a2, o2) {
            var s2 = c(a2, Ot);
            if (s2) {
              var l2 = ae(s2).sseEventSource, u2 = function(e2) {
                if (!Tt(s2)) {
                  if (!se(a2)) {
                    l2.removeEventListener(o2, u2);
                    return;
                  }
                  var t2 = e2.data;
                  R(a2, function(e3) {
                    t2 = e3.transformResponse(t2, null, a2);
                  });
                  var r2 = wr(a2), n2 = ye(a2), i2 = T(a2);
                  je(r2.swapStyle, n2, a2, t2, i2), nr(i2.tasks), ce(a2, "htmx:sseMessage", e2);
                }
              };
              ae(a2).sseListener = u2, l2.addEventListener(o2, u2);
            } else
              fe(a2, "htmx:noSSESourceError");
          }
          function Rt(e2, t2, r2) {
            var n2 = c(e2, Ot);
            if (n2) {
              var i2 = ae(n2).sseEventSource, a2 = function() {
                Tt(n2) || (se(e2) ? t2(e2) : i2.removeEventListener(r2, a2));
              };
              ae(e2).sseListener = a2, i2.addEventListener(r2, a2);
            } else
              fe(e2, "htmx:noSSESourceError");
          }
          function Tt(e2) {
            if (!se(e2))
              return ae(e2).sseEventSource.close(), !0;
          }
          function Ot(e2) {
            return ae(e2).sseEventSource != null;
          }
          function qt(e2, t2, r2, n2) {
            var i2 = function() {
              r2.loaded || (r2.loaded = !0, t2(e2));
            };
            n2 > 0 ? setTimeout(i2, n2) : i2();
          }
          function Ht(t2, i2, e2) {
            var a2 = !1;
            return oe(w, function(r2) {
              if (o(t2, "hx-" + r2)) {
                var n2 = te(t2, "hx-" + r2);
                a2 = !0, i2.path = n2, i2.verb = r2, e2.forEach(function(e3) {
                  Lt(t2, e3, i2, function(e4, t3) {
                    if (v(e4, Q.config.disableSelector)) {
                      p(e4);
                      return;
                    }
                    he(r2, n2, e4, t3);
                  });
                });
              }
            }), a2;
          }
          function Lt(n2, e2, t2, r2) {
            if (e2.sseEvent)
              Rt(n2, r2, e2.sseEvent);
            else if (e2.trigger === "revealed")
              gt(), ht(n2, r2, t2, e2), mt(n2);
            else if (e2.trigger === "intersect") {
              var i2 = {};
              e2.root && (i2.root = ue(n2, e2.root)), e2.threshold && (i2.threshold = parseFloat(e2.threshold));
              var a2 = new IntersectionObserver(function(e3) {
                for (var t3 = 0; t3 < e3.length; t3++) {
                  var r3 = e3[t3];
                  if (r3.isIntersecting) {
                    ce(n2, "intersect");
                    break;
                  }
                }
              }, i2);
              a2.observe(n2), ht(n2, r2, t2, e2);
            } else e2.trigger === "load" ? ct(e2, n2, Wt("load", { elt: n2 })) || qt(n2, r2, t2, e2.delay) : e2.pollInterval > 0 ? (t2.polling = !0, ot(n2, r2, e2)) : ht(n2, r2, t2, e2);
          }
          function At(e2) {
            if (Q.config.allowScriptTags && (e2.type === "text/javascript" || e2.type === "module" || e2.type === "")) {
              var t2 = re().createElement("script");
              oe(e2.attributes, function(e3) {
                t2.setAttribute(e3.name, e3.value);
              }), t2.textContent = e2.textContent, t2.async = !1, Q.config.inlineScriptNonce && (t2.nonce = Q.config.inlineScriptNonce);
              var r2 = e2.parentElement;
              try {
                r2.insertBefore(t2, e2);
              } catch (e3) {
                b(e3);
              } finally {
                e2.parentElement && e2.parentElement.removeChild(e2);
              }
            }
          }
          function Nt(e2) {
            h(e2, "script") && At(e2), oe(f(e2, "script"), function(e3) {
              At(e3);
            });
          }
          function It(e2) {
            for (var t2 = e2.attributes, r2 = 0; r2 < t2.length; r2++) {
              var n2 = t2[r2].name;
              if (s(n2, "hx-on:") || s(n2, "data-hx-on:") || s(n2, "hx-on-") || s(n2, "data-hx-on-"))
                return !0;
            }
            return !1;
          }
          function kt(e2) {
            var t2 = null, r2 = [];
            if (It(e2) && r2.push(e2), document.evaluate)
              for (var n2 = document.evaluate('.//*[@*[ starts-with(name(), "hx-on:") or starts-with(name(), "data-hx-on:") or starts-with(name(), "hx-on-") or starts-with(name(), "data-hx-on-") ]]', e2); t2 = n2.iterateNext(); ) r2.push(t2);
            else
              for (var i2 = e2.getElementsByTagName("*"), a2 = 0; a2 < i2.length; a2++)
                It(i2[a2]) && r2.push(i2[a2]);
            return r2;
          }
          function Pt(e2) {
            if (e2.querySelectorAll) {
              var t2 = ", [hx-boost] a, [data-hx-boost] a, a[hx-boost], a[data-hx-boost]", r2 = e2.querySelectorAll(i + t2 + ", form, [type='submit'], [hx-sse], [data-hx-sse], [hx-ws], [data-hx-ws], [hx-ext], [data-hx-ext], [hx-trigger], [data-hx-trigger], [hx-on], [data-hx-on]");
              return r2;
            } else
              return [];
          }
          function Mt(e2) {
            var t2 = v(e2.target, "button, input[type='submit']"), r2 = Dt(e2);
            r2 && (r2.lastButtonClicked = t2);
          }
          function Xt(e2) {
            var t2 = Dt(e2);
            t2 && (t2.lastButtonClicked = null);
          }
          function Dt(e2) {
            var t2 = v(e2.target, "button, input[type='submit']");
            if (t2) {
              var r2 = g("#" + ee(t2, "form")) || v(t2, "form");
              if (r2)
                return ae(r2);
            }
          }
          function Ut(e2) {
            e2.addEventListener("click", Mt), e2.addEventListener("focusin", Mt), e2.addEventListener("focusout", Xt);
          }
          function Bt(e2) {
            for (var t2 = Ye(e2), r2 = 0, n2 = 0; n2 < t2.length; n2++) {
              let i2 = t2[n2];
              i2 === "{" ? r2++ : i2 === "}" && r2--;
            }
            return r2;
          }
          function Ft(t2, e2, r2) {
            var n2 = ae(t2);
            Array.isArray(n2.onHandlers) || (n2.onHandlers = []);
            var i2, a2 = function(e3) {
              return Tr(t2, function() {
                i2 || (i2 = new Function("event", r2)), i2.call(t2, e3);
              });
            };
            t2.addEventListener(e2, a2), n2.onHandlers.push({ event: e2, listener: a2 });
          }
          function Vt(e2) {
            var t2 = te(e2, "hx-on");
            if (t2) {
              for (var r2 = {}, n2 = t2.split(`
`), i2 = null, a2 = 0; n2.length > 0; ) {
                var o2 = n2.shift(), s2 = o2.match(/^\s*([a-zA-Z:\-\.]+:)(.*)/);
                a2 === 0 && s2 ? (o2.split(":"), i2 = s2[1].slice(0, -1), r2[i2] = s2[2]) : r2[i2] += o2, a2 += Bt(o2);
              }
              for (var l2 in r2)
                Ft(e2, l2, r2[l2]);
            }
          }
          function jt(e2) {
            Ae(e2);
            for (var t2 = 0; t2 < e2.attributes.length; t2++) {
              var r2 = e2.attributes[t2].name, n2 = e2.attributes[t2].value;
              if (s(r2, "hx-on") || s(r2, "data-hx-on")) {
                var i2 = r2.indexOf("-on") + 3, a2 = r2.slice(i2, i2 + 1);
                if (a2 === "-" || a2 === ":") {
                  var o2 = r2.slice(i2 + 1);
                  s(o2, ":") ? o2 = "htmx" + o2 : s(o2, "-") ? o2 = "htmx:" + o2.slice(1) : s(o2, "htmx-") && (o2 = "htmx:" + o2.slice(5)), Ft(e2, o2, n2);
                }
              }
            }
          }
          function _t(t2) {
            if (v(t2, Q.config.disableSelector)) {
              p(t2);
              return;
            }
            var r2 = ae(t2);
            if (r2.initHash !== Le(t2)) {
              Ne(t2), r2.initHash = Le(t2), Vt(t2), ce(t2, "htmx:beforeProcessNode"), t2.value && (r2.lastValue = t2.value);
              var e2 = it(t2), n2 = Ht(t2, r2, e2);
              n2 || (ne(t2, "hx-boost") === "true" ? lt(t2, r2, e2) : o(t2, "hx-trigger") && e2.forEach(function(e3) {
                Lt(t2, e3, r2, function() {
                });
              })), (t2.tagName === "FORM" || ee(t2, "type") === "submit" && o(t2, "form")) && Ut(t2);
              var i2 = te(t2, "hx-sse");
              i2 && St(t2, r2, i2);
              var a2 = te(t2, "hx-ws");
              a2 && pt(t2, r2, a2), ce(t2, "htmx:afterProcessNode");
            }
          }
          function zt(e2) {
            if (e2 = g(e2), v(e2, Q.config.disableSelector)) {
              p(e2);
              return;
            }
            _t(e2), oe(Pt(e2), function(e3) {
              _t(e3);
            }), oe(kt(e2), jt);
          }
          function $t(e2) {
            return e2.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
          }
          function Wt(e2, t2) {
            var r2;
            return window.CustomEvent && typeof window.CustomEvent == "function" ? r2 = new CustomEvent(e2, { bubbles: !0, cancelable: !0, detail: t2 }) : (r2 = re().createEvent("CustomEvent"), r2.initCustomEvent(e2, !0, !0, t2)), r2;
          }
          function fe(e2, t2, r2) {
            ce(e2, t2, le({ error: t2 }, r2));
          }
          function Gt(e2) {
            return e2 === "htmx:afterProcessNode";
          }
          function R(e2, t2) {
            oe(Fr(e2), function(e3) {
              try {
                t2(e3);
              } catch (e4) {
                b(e4);
              }
            });
          }
          function b(e2) {
            console.error ? console.error(e2) : console.log && console.log("ERROR: ", e2);
          }
          function ce(e2, t2, r2) {
            e2 = g(e2), r2 == null && (r2 = {}), r2.elt = e2;
            var n2 = Wt(t2, r2);
            Q.logger && !Gt(t2) && Q.logger(e2, t2, r2), r2.error && (b(r2.error), ce(e2, "htmx:error", { errorInfo: r2 }));
            var i2 = e2.dispatchEvent(n2), a2 = $t(t2);
            if (i2 && a2 !== t2) {
              var o2 = Wt(a2, n2.detail);
              i2 = i2 && e2.dispatchEvent(o2);
            }
            return R(e2, function(e3) {
              i2 = i2 && e3.onEvent(t2, n2) !== !1 && !n2.defaultPrevented;
            }), i2;
          }
          var Jt = location.pathname + location.search;
          function Zt() {
            var e2 = re().querySelector("[hx-history-elt],[data-hx-history-elt]");
            return e2 || re().body;
          }
          function Kt(e2, t2, r2, n2) {
            if (U()) {
              if (Q.config.historyCacheSize <= 0) {
                localStorage.removeItem("htmx-history-cache");
                return;
              }
              e2 = B(e2);
              for (var i2 = E(localStorage.getItem("htmx-history-cache")) || [], a2 = 0; a2 < i2.length; a2++)
                if (i2[a2].url === e2) {
                  i2.splice(a2, 1);
                  break;
                }
              var o2 = { url: e2, content: t2, title: r2, scroll: n2 };
              for (ce(re().body, "htmx:historyItemCreated", { item: o2, cache: i2 }), i2.push(o2); i2.length > Q.config.historyCacheSize; )
                i2.shift();
              for (; i2.length > 0; )
                try {
                  localStorage.setItem("htmx-history-cache", JSON.stringify(i2));
                  break;
                } catch (e3) {
                  fe(re().body, "htmx:historyCacheError", { cause: e3, cache: i2 }), i2.shift();
                }
            }
          }
          function Yt(e2) {
            if (!U())
              return null;
            e2 = B(e2);
            for (var t2 = E(localStorage.getItem("htmx-history-cache")) || [], r2 = 0; r2 < t2.length; r2++)
              if (t2[r2].url === e2)
                return t2[r2];
            return null;
          }
          function Qt(e2) {
            var t2 = Q.config.requestClass, r2 = e2.cloneNode(!0);
            return oe(f(r2, "." + t2), function(e3) {
              n(e3, t2);
            }), r2.innerHTML;
          }
          function er() {
            var e2 = Zt(), t2 = Jt || location.pathname + location.search, r2;
            try {
              r2 = re().querySelector('[hx-history="false" i],[data-hx-history="false" i]');
            } catch {
              r2 = re().querySelector('[hx-history="false"],[data-hx-history="false"]');
            }
            r2 || (ce(re().body, "htmx:beforeHistorySave", { path: t2, historyElt: e2 }), Kt(t2, Qt(e2), re().title, window.scrollY)), Q.config.historyEnabled && history.replaceState({ htmx: !0 }, re().title, window.location.href);
          }
          function tr(e2) {
            Q.config.getCacheBusterParam && (e2 = e2.replace(/org\.htmx\.cache-buster=[^&]*&?/, ""), (G(e2, "&") || G(e2, "?")) && (e2 = e2.slice(0, -1))), Q.config.historyEnabled && history.pushState({ htmx: !0 }, "", e2), Jt = e2;
          }
          function rr(e2) {
            Q.config.historyEnabled && history.replaceState({ htmx: !0 }, "", e2), Jt = e2;
          }
          function nr(e2) {
            oe(e2, function(e3) {
              e3.call();
            });
          }
          function ir(a2) {
            var e2 = new XMLHttpRequest(), o2 = { path: a2, xhr: e2 };
            ce(re().body, "htmx:historyCacheMiss", o2), e2.open("GET", a2, !0), e2.setRequestHeader("HX-Request", "true"), e2.setRequestHeader("HX-History-Restore-Request", "true"), e2.setRequestHeader("HX-Current-URL", re().location.href), e2.onload = function() {
              if (this.status >= 200 && this.status < 400) {
                ce(re().body, "htmx:historyCacheMissLoad", o2);
                var e3 = l(this.response);
                e3 = e3.querySelector("[hx-history-elt],[data-hx-history-elt]") || e3;
                var t2 = Zt(), r2 = T(t2), n2 = Ve(this.response);
                if (n2) {
                  var i2 = C("title");
                  i2 ? i2.innerHTML = n2 : window.document.title = n2;
                }
                Ue(t2, e3, r2), nr(r2.tasks), Jt = a2, ce(re().body, "htmx:historyRestore", { path: a2, cacheMiss: !0, serverResponse: this.response });
              } else
                fe(re().body, "htmx:historyCacheMissLoadError", o2);
            }, e2.send();
          }
          function ar(e2) {
            er(), e2 = e2 || location.pathname + location.search;
            var t2 = Yt(e2);
            if (t2) {
              var r2 = l(t2.content), n2 = Zt(), i2 = T(n2);
              Ue(n2, r2, i2), nr(i2.tasks), document.title = t2.title, setTimeout(function() {
                window.scrollTo(0, t2.scroll);
              }, 0), Jt = e2, ce(re().body, "htmx:historyRestore", { path: e2, item: t2 });
            } else
              Q.config.refreshOnHistoryMiss ? window.location.reload(!0) : ir(e2);
          }
          function or(e2) {
            var t2 = pe(e2, "hx-indicator");
            return t2 == null && (t2 = [e2]), oe(t2, function(e3) {
              var t3 = ae(e3);
              t3.requestCount = (t3.requestCount || 0) + 1, e3.classList.add.call(e3.classList, Q.config.requestClass);
            }), t2;
          }
          function sr(e2) {
            var t2 = pe(e2, "hx-disabled-elt");
            return t2 == null && (t2 = []), oe(t2, function(e3) {
              var t3 = ae(e3);
              t3.requestCount = (t3.requestCount || 0) + 1, e3.setAttribute("disabled", "");
            }), t2;
          }
          function lr(e2, t2) {
            oe(e2, function(e3) {
              var t3 = ae(e3);
              t3.requestCount = (t3.requestCount || 0) - 1, t3.requestCount === 0 && e3.classList.remove.call(e3.classList, Q.config.requestClass);
            }), oe(t2, function(e3) {
              var t3 = ae(e3);
              t3.requestCount = (t3.requestCount || 0) - 1, t3.requestCount === 0 && e3.removeAttribute("disabled");
            });
          }
          function ur(e2, t2) {
            for (var r2 = 0; r2 < e2.length; r2++) {
              var n2 = e2[r2];
              if (n2.isSameNode(t2))
                return !0;
            }
            return !1;
          }
          function fr(e2) {
            return e2.name === "" || e2.name == null || e2.disabled || v(e2, "fieldset[disabled]") || e2.type === "button" || e2.type === "submit" || e2.tagName === "image" || e2.tagName === "reset" || e2.tagName === "file" ? !1 : e2.type === "checkbox" || e2.type === "radio" ? e2.checked : !0;
          }
          function cr(e2, t2, r2) {
            if (e2 != null && t2 != null) {
              var n2 = r2[e2];
              n2 === void 0 ? r2[e2] = t2 : Array.isArray(n2) ? Array.isArray(t2) ? r2[e2] = n2.concat(t2) : n2.push(t2) : Array.isArray(t2) ? r2[e2] = [n2].concat(t2) : r2[e2] = [n2, t2];
            }
          }
          function hr(t2, r2, n2, e2, i2) {
            if (!(e2 == null || ur(t2, e2))) {
              if (t2.push(e2), fr(e2)) {
                var a2 = ee(e2, "name"), o2 = e2.value;
                e2.multiple && e2.tagName === "SELECT" && (o2 = M(e2.querySelectorAll("option:checked")).map(function(e3) {
                  return e3.value;
                })), e2.files && (o2 = M(e2.files)), cr(a2, o2, r2), i2 && vr(e2, n2);
              }
              if (h(e2, "form")) {
                var s2 = e2.elements;
                oe(s2, function(e3) {
                  hr(t2, r2, n2, e3, i2);
                });
              }
            }
          }
          function vr(e2, t2) {
            e2.willValidate && (ce(e2, "htmx:validation:validate"), e2.checkValidity() || (t2.push({ elt: e2, message: e2.validationMessage, validity: e2.validity }), ce(e2, "htmx:validation:failed", { message: e2.validationMessage, validity: e2.validity })));
          }
          function dr(e2, t2) {
            var r2 = [], n2 = {}, i2 = {}, a2 = [], o2 = ae(e2);
            o2.lastButtonClicked && !se(o2.lastButtonClicked) && (o2.lastButtonClicked = null);
            var s2 = h(e2, "form") && e2.noValidate !== !0 || te(e2, "hx-validate") === "true";
            if (o2.lastButtonClicked && (s2 = s2 && o2.lastButtonClicked.formNoValidate !== !0), t2 !== "get" && hr(r2, i2, a2, v(e2, "form"), s2), hr(r2, n2, a2, e2, s2), o2.lastButtonClicked || e2.tagName === "BUTTON" || e2.tagName === "INPUT" && ee(e2, "type") === "submit") {
              var l2 = o2.lastButtonClicked || e2, u2 = ee(l2, "name");
              cr(u2, l2.value, i2);
            }
            var f2 = pe(e2, "hx-include");
            return oe(f2, function(e3) {
              hr(r2, n2, a2, e3, s2), h(e3, "form") || oe(e3.querySelectorAll(rt), function(e4) {
                hr(r2, n2, a2, e4, s2);
              });
            }), n2 = le(n2, i2), { errors: a2, values: n2 };
          }
          function gr(e2, t2, r2) {
            e2 !== "" && (e2 += "&"), String(r2) === "[object Object]" && (r2 = JSON.stringify(r2));
            var n2 = encodeURIComponent(r2);
            return e2 += encodeURIComponent(t2) + "=" + n2, e2;
          }
          function mr(e2) {
            var t2 = "";
            for (var r2 in e2)
              if (e2.hasOwnProperty(r2)) {
                var n2 = e2[r2];
                Array.isArray(n2) ? oe(n2, function(e3) {
                  t2 = gr(t2, r2, e3);
                }) : t2 = gr(t2, r2, n2);
              }
            return t2;
          }
          function pr(e2) {
            var t2 = new FormData();
            for (var r2 in e2)
              if (e2.hasOwnProperty(r2)) {
                var n2 = e2[r2];
                Array.isArray(n2) ? oe(n2, function(e3) {
                  t2.append(r2, e3);
                }) : t2.append(r2, n2);
              }
            return t2;
          }
          function xr(e2, t2, r2) {
            var n2 = { "HX-Request": "true", "HX-Trigger": ee(e2, "id"), "HX-Trigger-Name": ee(e2, "name"), "HX-Target": te(t2, "id"), "HX-Current-URL": re().location.href };
            return Rr(e2, "hx-headers", !1, n2), r2 !== void 0 && (n2["HX-Prompt"] = r2), ae(e2).boosted && (n2["HX-Boosted"] = "true"), n2;
          }
          function yr(t2, e2) {
            var r2 = ne(e2, "hx-params");
            if (r2) {
              if (r2 === "none")
                return {};
              if (r2 === "*")
                return t2;
              if (r2.indexOf("not ") === 0)
                return oe(r2.substr(4).split(","), function(e3) {
                  e3 = e3.trim(), delete t2[e3];
                }), t2;
              var n2 = {};
              return oe(r2.split(","), function(e3) {
                e3 = e3.trim(), n2[e3] = t2[e3];
              }), n2;
            } else
              return t2;
          }
          function br(e2) {
            return ee(e2, "href") && ee(e2, "href").indexOf("#") >= 0;
          }
          function wr(e2, t2) {
            var r2 = t2 || ne(e2, "hx-swap"), n2 = { swapStyle: ae(e2).boosted ? "innerHTML" : Q.config.defaultSwapStyle, swapDelay: Q.config.defaultSwapDelay, settleDelay: Q.config.defaultSettleDelay };
            if (Q.config.scrollIntoViewOnBoost && ae(e2).boosted && !br(e2) && (n2.show = "top"), r2) {
              var i2 = D(r2);
              if (i2.length > 0)
                for (var a2 = 0; a2 < i2.length; a2++) {
                  var o2 = i2[a2];
                  if (o2.indexOf("swap:") === 0)
                    n2.swapDelay = d(o2.substr(5));
                  else if (o2.indexOf("settle:") === 0)
                    n2.settleDelay = d(o2.substr(7));
                  else if (o2.indexOf("transition:") === 0)
                    n2.transition = o2.substr(11) === "true";
                  else if (o2.indexOf("ignoreTitle:") === 0)
                    n2.ignoreTitle = o2.substr(12) === "true";
                  else if (o2.indexOf("scroll:") === 0) {
                    var s2 = o2.substr(7), l2 = s2.split(":"), u2 = l2.pop(), f2 = l2.length > 0 ? l2.join(":") : null;
                    n2.scroll = u2, n2.scrollTarget = f2;
                  } else if (o2.indexOf("show:") === 0) {
                    var c2 = o2.substr(5), l2 = c2.split(":"), h2 = l2.pop(), f2 = l2.length > 0 ? l2.join(":") : null;
                    n2.show = h2, n2.showTarget = f2;
                  } else if (o2.indexOf("focus-scroll:") === 0) {
                    var v2 = o2.substr(13);
                    n2.focusScroll = v2 == "true";
                  } else a2 == 0 ? n2.swapStyle = o2 : b("Unknown modifier in hx-swap: " + o2);
                }
            }
            return n2;
          }
          function Sr(e2) {
            return ne(e2, "hx-encoding") === "multipart/form-data" || h(e2, "form") && ee(e2, "enctype") === "multipart/form-data";
          }
          function Er(t2, r2, n2) {
            var i2 = null;
            return R(r2, function(e2) {
              i2 == null && (i2 = e2.encodeParameters(t2, n2, r2));
            }), i2 ?? (Sr(r2) ? pr(n2) : mr(n2));
          }
          function T(e2) {
            return { tasks: [], elts: [e2] };
          }
          function Cr(e2, t2) {
            var r2 = e2[0], n2 = e2[e2.length - 1];
            if (t2.scroll) {
              var i2 = null;
              t2.scrollTarget && (i2 = ue(r2, t2.scrollTarget)), t2.scroll === "top" && (r2 || i2) && (i2 = i2 || r2, i2.scrollTop = 0), t2.scroll === "bottom" && (n2 || i2) && (i2 = i2 || n2, i2.scrollTop = i2.scrollHeight);
            }
            if (t2.show) {
              var i2 = null;
              if (t2.showTarget) {
                var a2 = t2.showTarget;
                t2.showTarget === "window" && (a2 = "body"), i2 = ue(r2, a2);
              }
              t2.show === "top" && (r2 || i2) && (i2 = i2 || r2, i2.scrollIntoView({ block: "start", behavior: Q.config.scrollBehavior })), t2.show === "bottom" && (n2 || i2) && (i2 = i2 || n2, i2.scrollIntoView({ block: "end", behavior: Q.config.scrollBehavior }));
            }
          }
          function Rr(e2, t2, r2, n2) {
            if (n2 == null && (n2 = {}), e2 == null)
              return n2;
            var i2 = te(e2, t2);
            if (i2) {
              var a2 = i2.trim(), o2 = r2;
              if (a2 === "unset")
                return null;
              a2.indexOf("javascript:") === 0 ? (a2 = a2.substr(11), o2 = !0) : a2.indexOf("js:") === 0 && (a2 = a2.substr(3), o2 = !0), a2.indexOf("{") !== 0 && (a2 = "{" + a2 + "}");
              var s2;
              o2 ? s2 = Tr(e2, function() {
                return Function("return (" + a2 + ")")();
              }, {}) : s2 = E(a2);
              for (var l2 in s2)
                s2.hasOwnProperty(l2) && n2[l2] == null && (n2[l2] = s2[l2]);
            }
            return Rr(u(e2), t2, r2, n2);
          }
          function Tr(e2, t2, r2) {
            return Q.config.allowEval ? t2() : (fe(e2, "htmx:evalDisallowedError"), r2);
          }
          function Or(e2, t2) {
            return Rr(e2, "hx-vars", !0, t2);
          }
          function qr(e2, t2) {
            return Rr(e2, "hx-vals", !1, t2);
          }
          function Hr(e2) {
            return le(Or(e2), qr(e2));
          }
          function Lr(t2, r2, n2) {
            if (n2 !== null)
              try {
                t2.setRequestHeader(r2, n2);
              } catch {
                t2.setRequestHeader(r2, encodeURIComponent(n2)), t2.setRequestHeader(r2 + "-URI-AutoEncoded", "true");
              }
          }
          function Ar(t2) {
            if (t2.responseURL && typeof URL < "u")
              try {
                var e2 = new URL(t2.responseURL);
                return e2.pathname + e2.search;
              } catch {
                fe(re().body, "htmx:badResponseUrl", { url: t2.responseURL });
              }
          }
          function O(e2, t2) {
            return t2.test(e2.getAllResponseHeaders());
          }
          function Nr(e2, t2, r2) {
            return e2 = e2.toLowerCase(), r2 ? r2 instanceof Element || I(r2, "String") ? he(e2, t2, null, null, { targetOverride: g(r2), returnPromise: !0 }) : he(e2, t2, g(r2.source), r2.event, { handler: r2.handler, headers: r2.headers, values: r2.values, targetOverride: g(r2.target), swapOverride: r2.swap, select: r2.select, returnPromise: !0 }) : he(e2, t2, null, null, { returnPromise: !0 });
          }
          function Ir(e2) {
            for (var t2 = []; e2; )
              t2.push(e2), e2 = e2.parentElement;
            return t2;
          }
          function kr(e2, t2, r2) {
            var n2, i2;
            if (typeof URL == "function") {
              i2 = new URL(t2, document.location.href);
              var a2 = document.location.origin;
              n2 = a2 === i2.origin;
            } else
              i2 = t2, n2 = s(t2, document.location.origin);
            return Q.config.selfRequestsOnly && !n2 ? !1 : ce(e2, "htmx:validateUrl", le({ url: i2, sameHost: n2 }, r2));
          }
          function he(t2, r2, n2, i2, a2, e2) {
            var o2 = null, s2 = null;
            if (a2 = a2 ?? {}, a2.returnPromise && typeof Promise < "u")
              var l2 = new Promise(function(e3, t3) {
                o2 = e3, s2 = t3;
              });
            n2 == null && (n2 = re().body);
            var M2 = a2.handler || Mr, X2 = a2.select || null;
            if (!se(n2))
              return ie(o2), l2;
            var u2 = a2.targetOverride || ye(n2);
            if (u2 == null || u2 == me)
              return fe(n2, "htmx:targetError", { target: te(n2, "hx-target") }), ie(s2), l2;
            var f2 = ae(n2), c2 = f2.lastButtonClicked;
            if (c2) {
              var h2 = ee(c2, "formaction");
              h2 != null && (r2 = h2);
              var v2 = ee(c2, "formmethod");
              v2 != null && v2.toLowerCase() !== "dialog" && (t2 = v2);
            }
            var d2 = ne(n2, "hx-confirm");
            if (e2 === void 0) {
              var D2 = function(e3) {
                return he(t2, r2, n2, i2, a2, !!e3);
              }, U2 = { target: u2, elt: n2, path: r2, verb: t2, triggeringEvent: i2, etc: a2, issueRequest: D2, question: d2 };
              if (ce(n2, "htmx:confirm", U2) === !1)
                return ie(o2), l2;
            }
            var g2 = n2, m2 = ne(n2, "hx-sync"), p2 = null, x2 = !1;
            if (m2) {
              var B2 = m2.split(":"), F2 = B2[0].trim();
              if (F2 === "this" ? g2 = xe(n2, "hx-sync") : g2 = ue(n2, F2), m2 = (B2[1] || "drop").trim(), f2 = ae(g2), m2 === "drop" && f2.xhr && f2.abortable !== !0)
                return ie(o2), l2;
              if (m2 === "abort") {
                if (f2.xhr)
                  return ie(o2), l2;
                x2 = !0;
              } else if (m2 === "replace")
                ce(g2, "htmx:abort");
              else if (m2.indexOf("queue") === 0) {
                var V2 = m2.split(" ");
                p2 = (V2[1] || "last").trim();
              }
            }
            if (f2.xhr)
              if (f2.abortable)
                ce(g2, "htmx:abort");
              else {
                if (p2 == null) {
                  if (i2) {
                    var y2 = ae(i2);
                    y2 && y2.triggerSpec && y2.triggerSpec.queue && (p2 = y2.triggerSpec.queue);
                  }
                  p2 == null && (p2 = "last");
                }
                return f2.queuedRequests == null && (f2.queuedRequests = []), p2 === "first" && f2.queuedRequests.length === 0 ? f2.queuedRequests.push(function() {
                  he(t2, r2, n2, i2, a2);
                }) : p2 === "all" ? f2.queuedRequests.push(function() {
                  he(t2, r2, n2, i2, a2);
                }) : p2 === "last" && (f2.queuedRequests = [], f2.queuedRequests.push(function() {
                  he(t2, r2, n2, i2, a2);
                })), ie(o2), l2;
              }
            var b2 = new XMLHttpRequest();
            f2.xhr = b2, f2.abortable = x2;
            var w2 = function() {
              if (f2.xhr = null, f2.abortable = !1, f2.queuedRequests != null && f2.queuedRequests.length > 0) {
                var e3 = f2.queuedRequests.shift();
                e3();
              }
            }, j2 = ne(n2, "hx-prompt");
            if (j2) {
              var S2 = prompt(j2);
              if (S2 === null || !ce(n2, "htmx:prompt", { prompt: S2, target: u2 }))
                return ie(o2), w2(), l2;
            }
            if (d2 && !e2 && !confirm(d2))
              return ie(o2), w2(), l2;
            var E2 = xr(n2, u2, S2);
            t2 !== "get" && !Sr(n2) && (E2["Content-Type"] = "application/x-www-form-urlencoded"), a2.headers && (E2 = le(E2, a2.headers));
            var _2 = dr(n2, t2), C2 = _2.errors, R2 = _2.values;
            a2.values && (R2 = le(R2, a2.values));
            var z2 = Hr(n2), $2 = le(R2, z2), T2 = yr($2, n2);
            Q.config.getCacheBusterParam && t2 === "get" && (T2["org.htmx.cache-buster"] = ee(u2, "id") || "true"), (r2 == null || r2 === "") && (r2 = re().location.href);
            var O2 = Rr(n2, "hx-request"), W2 = ae(n2).boosted, q2 = Q.config.methodsThatUseUrlParams.indexOf(t2) >= 0, H2 = { boosted: W2, useUrlParams: q2, parameters: T2, unfilteredParameters: $2, headers: E2, target: u2, verb: t2, errors: C2, withCredentials: a2.credentials || O2.credentials || Q.config.withCredentials, timeout: a2.timeout || O2.timeout || Q.config.timeout, path: r2, triggeringEvent: i2 };
            if (!ce(n2, "htmx:configRequest", H2))
              return ie(o2), w2(), l2;
            if (r2 = H2.path, t2 = H2.verb, E2 = H2.headers, T2 = H2.parameters, C2 = H2.errors, q2 = H2.useUrlParams, C2 && C2.length > 0)
              return ce(n2, "htmx:validation:halted", H2), ie(o2), w2(), l2;
            var G2 = r2.split("#"), J2 = G2[0], L2 = G2[1], A2 = r2;
            if (q2) {
              A2 = J2;
              var Z2 = Object.keys(T2).length !== 0;
              Z2 && (A2.indexOf("?") < 0 ? A2 += "?" : A2 += "&", A2 += mr(T2), L2 && (A2 += "#" + L2));
            }
            if (!kr(n2, A2, H2))
              return fe(n2, "htmx:invalidPath", H2), ie(s2), l2;
            if (b2.open(t2.toUpperCase(), A2, !0), b2.overrideMimeType("text/html"), b2.withCredentials = H2.withCredentials, b2.timeout = H2.timeout, !O2.noHeaders) {
              for (var N2 in E2)
                if (E2.hasOwnProperty(N2)) {
                  var K2 = E2[N2];
                  Lr(b2, N2, K2);
                }
            }
            var I2 = { xhr: b2, target: u2, requestConfig: H2, etc: a2, boosted: W2, select: X2, pathInfo: { requestPath: r2, finalRequestPath: A2, anchor: L2 } };
            if (b2.onload = function() {
              try {
                var e3 = Ir(n2);
                if (I2.pathInfo.responsePath = Ar(b2), M2(n2, I2), lr(k2, P2), ce(n2, "htmx:afterRequest", I2), ce(n2, "htmx:afterOnLoad", I2), !se(n2)) {
                  for (var t3 = null; e3.length > 0 && t3 == null; ) {
                    var r3 = e3.shift();
                    se(r3) && (t3 = r3);
                  }
                  t3 && (ce(t3, "htmx:afterRequest", I2), ce(t3, "htmx:afterOnLoad", I2));
                }
                ie(o2), w2();
              } catch (e4) {
                throw fe(n2, "htmx:onLoadError", le({ error: e4 }, I2)), e4;
              }
            }, b2.onerror = function() {
              lr(k2, P2), fe(n2, "htmx:afterRequest", I2), fe(n2, "htmx:sendError", I2), ie(s2), w2();
            }, b2.onabort = function() {
              lr(k2, P2), fe(n2, "htmx:afterRequest", I2), fe(n2, "htmx:sendAbort", I2), ie(s2), w2();
            }, b2.ontimeout = function() {
              lr(k2, P2), fe(n2, "htmx:afterRequest", I2), fe(n2, "htmx:timeout", I2), ie(s2), w2();
            }, !ce(n2, "htmx:beforeRequest", I2))
              return ie(o2), w2(), l2;
            var k2 = or(n2), P2 = sr(n2);
            oe(["loadstart", "loadend", "progress", "abort"], function(t3) {
              oe([b2, b2.upload], function(e3) {
                e3.addEventListener(t3, function(e4) {
                  ce(n2, "htmx:xhr:" + t3, { lengthComputable: e4.lengthComputable, loaded: e4.loaded, total: e4.total });
                });
              });
            }), ce(n2, "htmx:beforeSend", I2);
            var Y2 = q2 ? null : Er(b2, n2, T2);
            return b2.send(Y2), l2;
          }
          function Pr(e2, t2) {
            var r2 = t2.xhr, n2 = null, i2 = null;
            if (O(r2, /HX-Push:/i) ? (n2 = r2.getResponseHeader("HX-Push"), i2 = "push") : O(r2, /HX-Push-Url:/i) ? (n2 = r2.getResponseHeader("HX-Push-Url"), i2 = "push") : O(r2, /HX-Replace-Url:/i) && (n2 = r2.getResponseHeader("HX-Replace-Url"), i2 = "replace"), n2)
              return n2 === "false" ? {} : { type: i2, path: n2 };
            var a2 = t2.pathInfo.finalRequestPath, o2 = t2.pathInfo.responsePath, s2 = ne(e2, "hx-push-url"), l2 = ne(e2, "hx-replace-url"), u2 = ae(e2).boosted, f2 = null, c2 = null;
            return s2 ? (f2 = "push", c2 = s2) : l2 ? (f2 = "replace", c2 = l2) : u2 && (f2 = "push", c2 = o2 || a2), c2 ? c2 === "false" ? {} : (c2 === "true" && (c2 = o2 || a2), t2.pathInfo.anchor && c2.indexOf("#") === -1 && (c2 = c2 + "#" + t2.pathInfo.anchor), { type: f2, path: c2 }) : {};
          }
          function Mr(l2, u2) {
            var f2 = u2.xhr, c2 = u2.target, e2 = u2.etc, t2 = u2.requestConfig, h2 = u2.select;
            if (ce(l2, "htmx:beforeOnLoad", u2)) {
              if (O(f2, /HX-Trigger:/i) && _e(f2, "HX-Trigger", l2), O(f2, /HX-Location:/i)) {
                er();
                var r2 = f2.getResponseHeader("HX-Location"), v2;
                r2.indexOf("{") === 0 && (v2 = E(r2), r2 = v2.path, delete v2.path), Nr("GET", r2, v2).then(function() {
                  tr(r2);
                });
                return;
              }
              var n2 = O(f2, /HX-Refresh:/i) && f2.getResponseHeader("HX-Refresh") === "true";
              if (O(f2, /HX-Redirect:/i)) {
                location.href = f2.getResponseHeader("HX-Redirect"), n2 && location.reload();
                return;
              }
              if (n2) {
                location.reload();
                return;
              }
              O(f2, /HX-Retarget:/i) && (f2.getResponseHeader("HX-Retarget") === "this" ? u2.target = l2 : u2.target = ue(l2, f2.getResponseHeader("HX-Retarget")));
              var d2 = Pr(l2, u2), i2 = f2.status >= 200 && f2.status < 400 && f2.status !== 204, g2 = f2.response, a2 = f2.status >= 400, m2 = Q.config.ignoreTitle, o2 = le({ shouldSwap: i2, serverResponse: g2, isError: a2, ignoreTitle: m2 }, u2);
              if (ce(c2, "htmx:beforeSwap", o2)) {
                if (c2 = o2.target, g2 = o2.serverResponse, a2 = o2.isError, m2 = o2.ignoreTitle, u2.target = c2, u2.failed = a2, u2.successful = !a2, o2.shouldSwap) {
                  f2.status === 286 && at(l2), R(l2, function(e3) {
                    g2 = e3.transformResponse(g2, f2, l2);
                  }), d2.type && er();
                  var s2 = e2.swapOverride;
                  O(f2, /HX-Reswap:/i) && (s2 = f2.getResponseHeader("HX-Reswap"));
                  var v2 = wr(l2, s2);
                  v2.hasOwnProperty("ignoreTitle") && (m2 = v2.ignoreTitle), c2.classList.add(Q.config.swappingClass);
                  var p2 = null, x2 = null, y2 = function() {
                    try {
                      var e3 = document.activeElement, t3 = {};
                      try {
                        t3 = { elt: e3, start: e3 ? e3.selectionStart : null, end: e3 ? e3.selectionEnd : null };
                      } catch {
                      }
                      var r3;
                      h2 && (r3 = h2), O(f2, /HX-Reselect:/i) && (r3 = f2.getResponseHeader("HX-Reselect")), d2.type && (ce(re().body, "htmx:beforeHistoryUpdate", le({ history: d2 }, u2)), d2.type === "push" ? (tr(d2.path), ce(re().body, "htmx:pushedIntoHistory", { path: d2.path })) : (rr(d2.path), ce(re().body, "htmx:replacedInHistory", { path: d2.path })));
                      var n3 = T(c2);
                      if (je(v2.swapStyle, c2, l2, g2, n3, r3), t3.elt && !se(t3.elt) && ee(t3.elt, "id")) {
                        var i3 = document.getElementById(ee(t3.elt, "id")), a3 = { preventScroll: v2.focusScroll !== void 0 ? !v2.focusScroll : !Q.config.defaultFocusScroll };
                        if (i3) {
                          if (t3.start && i3.setSelectionRange)
                            try {
                              i3.setSelectionRange(t3.start, t3.end);
                            } catch {
                            }
                          i3.focus(a3);
                        }
                      }
                      if (c2.classList.remove(Q.config.swappingClass), oe(n3.elts, function(e4) {
                        e4.classList && e4.classList.add(Q.config.settlingClass), ce(e4, "htmx:afterSwap", u2);
                      }), O(f2, /HX-Trigger-After-Swap:/i)) {
                        var o3 = l2;
                        se(l2) || (o3 = re().body), _e(f2, "HX-Trigger-After-Swap", o3);
                      }
                      var s3 = function() {
                        if (oe(n3.tasks, function(e5) {
                          e5.call();
                        }), oe(n3.elts, function(e5) {
                          e5.classList && e5.classList.remove(Q.config.settlingClass), ce(e5, "htmx:afterSettle", u2);
                        }), u2.pathInfo.anchor) {
                          var e4 = re().getElementById(u2.pathInfo.anchor);
                          e4 && e4.scrollIntoView({ block: "start", behavior: "auto" });
                        }
                        if (n3.title && !m2) {
                          var t4 = C("title");
                          t4 ? t4.innerHTML = n3.title : window.document.title = n3.title;
                        }
                        if (Cr(n3.elts, v2), O(f2, /HX-Trigger-After-Settle:/i)) {
                          var r4 = l2;
                          se(l2) || (r4 = re().body), _e(f2, "HX-Trigger-After-Settle", r4);
                        }
                        ie(p2);
                      };
                      v2.settleDelay > 0 ? setTimeout(s3, v2.settleDelay) : s3();
                    } catch (e4) {
                      throw fe(l2, "htmx:swapError", u2), ie(x2), e4;
                    }
                  }, b2 = Q.config.globalViewTransitions;
                  if (v2.hasOwnProperty("transition") && (b2 = v2.transition), b2 && ce(l2, "htmx:beforeTransition", u2) && typeof Promise < "u" && document.startViewTransition) {
                    var w2 = new Promise(function(e3, t3) {
                      p2 = e3, x2 = t3;
                    }), S2 = y2;
                    y2 = function() {
                      document.startViewTransition(function() {
                        return S2(), w2;
                      });
                    };
                  }
                  v2.swapDelay > 0 ? setTimeout(y2, v2.swapDelay) : y2();
                }
                a2 && fe(l2, "htmx:responseError", le({ error: "Response Status Error Code " + f2.status + " from " + u2.pathInfo.requestPath }, u2));
              }
            }
          }
          var Xr = {};
          function Dr() {
            return { init: function(e2) {
              return null;
            }, onEvent: function(e2, t2) {
              return !0;
            }, transformResponse: function(e2, t2, r2) {
              return e2;
            }, isInlineSwap: function(e2) {
              return !1;
            }, handleSwap: function(e2, t2, r2, n2) {
              return !1;
            }, encodeParameters: function(e2, t2, r2) {
              return null;
            } };
          }
          function Ur(e2, t2) {
            t2.init && t2.init(r), Xr[e2] = le(Dr(), t2);
          }
          function Br(e2) {
            delete Xr[e2];
          }
          function Fr(e2, r2, n2) {
            if (e2 == null)
              return r2;
            r2 == null && (r2 = []), n2 == null && (n2 = []);
            var t2 = te(e2, "hx-ext");
            return t2 && oe(t2.split(","), function(e3) {
              if (e3 = e3.replace(/ /g, ""), e3.slice(0, 7) == "ignore:") {
                n2.push(e3.slice(7));
                return;
              }
              if (n2.indexOf(e3) < 0) {
                var t3 = Xr[e3];
                t3 && r2.indexOf(t3) < 0 && r2.push(t3);
              }
            }), Fr(u(e2), r2, n2);
          }
          var Vr = !1;
          re().addEventListener("DOMContentLoaded", function() {
            Vr = !0;
          });
          function jr(e2) {
            Vr || re().readyState === "complete" ? e2() : re().addEventListener("DOMContentLoaded", e2);
          }
          function _r() {
            Q.config.includeIndicatorStyles !== !1 && re().head.insertAdjacentHTML("beforeend", "<style>                      ." + Q.config.indicatorClass + "{opacity:0}                      ." + Q.config.requestClass + " ." + Q.config.indicatorClass + "{opacity:1; transition: opacity 200ms ease-in;}                      ." + Q.config.requestClass + "." + Q.config.indicatorClass + "{opacity:1; transition: opacity 200ms ease-in;}                    </style>");
          }
          function zr() {
            var e2 = re().querySelector('meta[name="htmx-config"]');
            return e2 ? E(e2.content) : null;
          }
          function $r() {
            var e2 = zr();
            e2 && (Q.config = le(Q.config, e2));
          }
          return jr(function() {
            $r(), _r();
            var e2 = re().body;
            zt(e2);
            var t2 = re().querySelectorAll("[hx-trigger='restored'],[data-hx-trigger='restored']");
            e2.addEventListener("htmx:abort", function(e3) {
              var t3 = e3.target, r3 = ae(t3);
              r3 && r3.xhr && r3.xhr.abort();
            });
            let r2 = window.onpopstate ? window.onpopstate.bind(window) : null;
            window.onpopstate = function(e3) {
              e3.state && e3.state.htmx ? (ar(), oe(t2, function(e4) {
                ce(e4, "htmx:restored", { document: re(), triggerEvent: ce });
              })) : r2 && r2(e3);
            }, setTimeout(function() {
              ce(e2, "htmx:load", {}), e2 = null;
            }, 0);
          }), Q;
        }();
      });
    }
  });

  // pkg/admin/src/main.ts
  var import_htmx = __toESM(require_htmx_min(), 1);

  // pkg/admin/src/wasm_exec.js
  (() => {
    let enosys = () => {
      let err = new Error("not implemented");
      return err.code = "ENOSYS", err;
    };
    if (!globalThis.fs) {
      let outputBuf = "";
      globalThis.fs = {
        constants: { O_WRONLY: -1, O_RDWR: -1, O_CREAT: -1, O_TRUNC: -1, O_APPEND: -1, O_EXCL: -1 },
        // unused
        writeSync(fd, buf) {
          outputBuf += decoder.decode(buf);
          let nl = outputBuf.lastIndexOf(`
`);
          return nl != -1 && (console.log(outputBuf.substring(0, nl)), outputBuf = outputBuf.substring(nl + 1)), buf.length;
        },
        write(fd, buf, offset, length, position, callback) {
          if (offset !== 0 || length !== buf.length || position !== null) {
            callback(enosys());
            return;
          }
          let n2 = this.writeSync(fd, buf);
          callback(null, n2);
        },
        chmod(path, mode, callback) {
          callback(enosys());
        },
        chown(path, uid, gid, callback) {
          callback(enosys());
        },
        close(fd, callback) {
          callback(enosys());
        },
        fchmod(fd, mode, callback) {
          callback(enosys());
        },
        fchown(fd, uid, gid, callback) {
          callback(enosys());
        },
        fstat(fd, callback) {
          callback(enosys());
        },
        fsync(fd, callback) {
          callback(null);
        },
        ftruncate(fd, length, callback) {
          callback(enosys());
        },
        lchown(path, uid, gid, callback) {
          callback(enosys());
        },
        link(path, link, callback) {
          callback(enosys());
        },
        lstat(path, callback) {
          callback(enosys());
        },
        mkdir(path, perm, callback) {
          callback(enosys());
        },
        open(path, flags, mode, callback) {
          callback(enosys());
        },
        read(fd, buffer, offset, length, position, callback) {
          callback(enosys());
        },
        readdir(path, callback) {
          callback(enosys());
        },
        readlink(path, callback) {
          callback(enosys());
        },
        rename(from, to, callback) {
          callback(enosys());
        },
        rmdir(path, callback) {
          callback(enosys());
        },
        stat(path, callback) {
          callback(enosys());
        },
        symlink(path, link, callback) {
          callback(enosys());
        },
        truncate(path, length, callback) {
          callback(enosys());
        },
        unlink(path, callback) {
          callback(enosys());
        },
        utimes(path, atime, mtime, callback) {
          callback(enosys());
        }
      };
    }
    if (!globalThis.crypto)
      throw new Error("globalThis.crypto is not available, polyfill required (crypto.getRandomValues only)");
    if (!globalThis.performance)
      throw new Error("globalThis.performance is not available, polyfill required (performance.now only)");
    if (!globalThis.TextEncoder)
      throw new Error("globalThis.TextEncoder is not available, polyfill required");
    if (!globalThis.TextDecoder)
      throw new Error("globalThis.TextDecoder is not available, polyfill required");
    let encoder = new TextEncoder("utf-8"), decoder = new TextDecoder("utf-8");
    globalThis.Go = class {
      constructor() {
        this.argv = ["js"], this.env = {}, this.exit = (code) => {
          code !== 0 && console.warn("exit code:", code);
        }, this._exitPromise = new Promise((resolve) => {
          this._resolveExitPromise = resolve;
        }), this._pendingEvent = null, this._scheduledTimeouts = /* @__PURE__ */ new Map(), this._nextCallbackTimeoutID = 1;
        let setInt64 = (addr, v2) => {
          this.mem.setUint32(addr + 0, v2, !0), this.mem.setUint32(addr + 4, Math.floor(v2 / 4294967296), !0);
        }, setInt32 = (addr, v2) => {
          this.mem.setUint32(addr + 0, v2, !0);
        }, getInt64 = (addr) => {
          let low = this.mem.getUint32(addr + 0, !0), high = this.mem.getInt32(addr + 4, !0);
          return low + high * 4294967296;
        }, loadValue = (addr) => {
          let f2 = this.mem.getFloat64(addr, !0);
          if (f2 === 0)
            return;
          if (!isNaN(f2))
            return f2;
          let id = this.mem.getUint32(addr, !0);
          return this._values[id];
        }, storeValue = (addr, v2) => {
          if (typeof v2 == "number" && v2 !== 0) {
            if (isNaN(v2)) {
              this.mem.setUint32(addr + 4, 2146959360, !0), this.mem.setUint32(addr, 0, !0);
              return;
            }
            this.mem.setFloat64(addr, v2, !0);
            return;
          }
          if (v2 === void 0) {
            this.mem.setFloat64(addr, 0, !0);
            return;
          }
          let id = this._ids.get(v2);
          id === void 0 && (id = this._idPool.pop(), id === void 0 && (id = this._values.length), this._values[id] = v2, this._goRefCounts[id] = 0, this._ids.set(v2, id)), this._goRefCounts[id]++;
          let typeFlag = 0;
          switch (typeof v2) {
            case "object":
              v2 !== null && (typeFlag = 1);
              break;
            case "string":
              typeFlag = 2;
              break;
            case "symbol":
              typeFlag = 3;
              break;
            case "function":
              typeFlag = 4;
              break;
          }
          this.mem.setUint32(addr + 4, 2146959360 | typeFlag, !0), this.mem.setUint32(addr, id, !0);
        }, loadSlice = (addr) => {
          let array = getInt64(addr + 0), len = getInt64(addr + 8);
          return new Uint8Array(this._inst.exports.mem.buffer, array, len);
        }, loadSliceOfValues = (addr) => {
          let array = getInt64(addr + 0), len = getInt64(addr + 8), a2 = new Array(len);
          for (let i2 = 0; i2 < len; i2++)
            a2[i2] = loadValue(array + i2 * 8);
          return a2;
        }, loadString = (addr) => {
          let saddr = getInt64(addr + 0), len = getInt64(addr + 8);
          return decoder.decode(new DataView(this._inst.exports.mem.buffer, saddr, len));
        }, timeOrigin = Date.now() - performance.now();
        this.importObject = {
          _gotest: {
            add: (a2, b2) => a2 + b2
          },
          gojs: {
            // Go's SP does not change as long as no Go code is running. Some operations (e.g. calls, getters and setters)
            // may synchronously trigger a Go event handler. This makes Go code get executed in the middle of the imported
            // function. A goroutine can switch to a new stack if the current stack is too small (see morestack function).
            // This changes the SP, thus we have to update the SP used by the imported function.
            // func wasmExit(code int32)
            "runtime.wasmExit": (sp) => {
              sp >>>= 0;
              let code = this.mem.getInt32(sp + 8, !0);
              this.exited = !0, delete this._inst, delete this._values, delete this._goRefCounts, delete this._ids, delete this._idPool, this.exit(code);
            },
            // func wasmWrite(fd uintptr, p unsafe.Pointer, n int32)
            "runtime.wasmWrite": (sp) => {
              sp >>>= 0;
              let fd = getInt64(sp + 8), p2 = getInt64(sp + 16), n2 = this.mem.getInt32(sp + 24, !0);
              fs.writeSync(fd, new Uint8Array(this._inst.exports.mem.buffer, p2, n2));
            },
            // func resetMemoryDataView()
            "runtime.resetMemoryDataView": (sp) => {
              sp >>>= 0, this.mem = new DataView(this._inst.exports.mem.buffer);
            },
            // func nanotime1() int64
            "runtime.nanotime1": (sp) => {
              sp >>>= 0, setInt64(sp + 8, (timeOrigin + performance.now()) * 1e6);
            },
            // func walltime() (sec int64, nsec int32)
            "runtime.walltime": (sp) => {
              sp >>>= 0;
              let msec = (/* @__PURE__ */ new Date()).getTime();
              setInt64(sp + 8, msec / 1e3), this.mem.setInt32(sp + 16, msec % 1e3 * 1e6, !0);
            },
            // func scheduleTimeoutEvent(delay int64) int32
            "runtime.scheduleTimeoutEvent": (sp) => {
              sp >>>= 0;
              let id = this._nextCallbackTimeoutID;
              this._nextCallbackTimeoutID++, this._scheduledTimeouts.set(id, setTimeout(
                () => {
                  for (this._resume(); this._scheduledTimeouts.has(id); )
                    console.warn("scheduleTimeoutEvent: missed timeout event"), this._resume();
                },
                getInt64(sp + 8)
              )), this.mem.setInt32(sp + 16, id, !0);
            },
            // func clearTimeoutEvent(id int32)
            "runtime.clearTimeoutEvent": (sp) => {
              sp >>>= 0;
              let id = this.mem.getInt32(sp + 8, !0);
              clearTimeout(this._scheduledTimeouts.get(id)), this._scheduledTimeouts.delete(id);
            },
            // func getRandomData(r []byte)
            "runtime.getRandomData": (sp) => {
              sp >>>= 0, crypto.getRandomValues(loadSlice(sp + 8));
            },
            // func finalizeRef(v ref)
            "syscall/js.finalizeRef": (sp) => {
              sp >>>= 0;
              let id = this.mem.getUint32(sp + 8, !0);
              if (this._goRefCounts[id]--, this._goRefCounts[id] === 0) {
                let v2 = this._values[id];
                this._values[id] = null, this._ids.delete(v2), this._idPool.push(id);
              }
            },
            // func stringVal(value string) ref
            "syscall/js.stringVal": (sp) => {
              sp >>>= 0, storeValue(sp + 24, loadString(sp + 8));
            },
            // func valueGet(v ref, p string) ref
            "syscall/js.valueGet": (sp) => {
              sp >>>= 0;
              let result = Reflect.get(loadValue(sp + 8), loadString(sp + 16));
              sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 32, result);
            },
            // func valueSet(v ref, p string, x ref)
            "syscall/js.valueSet": (sp) => {
              sp >>>= 0, Reflect.set(loadValue(sp + 8), loadString(sp + 16), loadValue(sp + 32));
            },
            // func valueDelete(v ref, p string)
            "syscall/js.valueDelete": (sp) => {
              sp >>>= 0, Reflect.deleteProperty(loadValue(sp + 8), loadString(sp + 16));
            },
            // func valueIndex(v ref, i int) ref
            "syscall/js.valueIndex": (sp) => {
              sp >>>= 0, storeValue(sp + 24, Reflect.get(loadValue(sp + 8), getInt64(sp + 16)));
            },
            // valueSetIndex(v ref, i int, x ref)
            "syscall/js.valueSetIndex": (sp) => {
              sp >>>= 0, Reflect.set(loadValue(sp + 8), getInt64(sp + 16), loadValue(sp + 24));
            },
            // func valueCall(v ref, m string, args []ref) (ref, bool)
            "syscall/js.valueCall": (sp) => {
              sp >>>= 0;
              try {
                let v2 = loadValue(sp + 8), m2 = Reflect.get(v2, loadString(sp + 16)), args = loadSliceOfValues(sp + 32), result = Reflect.apply(m2, v2, args);
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 56, result), this.mem.setUint8(sp + 64, 1);
              } catch (err) {
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 56, err), this.mem.setUint8(sp + 64, 0);
              }
            },
            // func valueInvoke(v ref, args []ref) (ref, bool)
            "syscall/js.valueInvoke": (sp) => {
              sp >>>= 0;
              try {
                let v2 = loadValue(sp + 8), args = loadSliceOfValues(sp + 16), result = Reflect.apply(v2, void 0, args);
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 40, result), this.mem.setUint8(sp + 48, 1);
              } catch (err) {
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 40, err), this.mem.setUint8(sp + 48, 0);
              }
            },
            // func valueNew(v ref, args []ref) (ref, bool)
            "syscall/js.valueNew": (sp) => {
              sp >>>= 0;
              try {
                let v2 = loadValue(sp + 8), args = loadSliceOfValues(sp + 16), result = Reflect.construct(v2, args);
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 40, result), this.mem.setUint8(sp + 48, 1);
              } catch (err) {
                sp = this._inst.exports.getsp() >>> 0, storeValue(sp + 40, err), this.mem.setUint8(sp + 48, 0);
              }
            },
            // func valueLength(v ref) int
            "syscall/js.valueLength": (sp) => {
              sp >>>= 0, setInt64(sp + 16, parseInt(loadValue(sp + 8).length));
            },
            // valuePrepareString(v ref) (ref, int)
            "syscall/js.valuePrepareString": (sp) => {
              sp >>>= 0;
              let str = encoder.encode(String(loadValue(sp + 8)));
              storeValue(sp + 16, str), setInt64(sp + 24, str.length);
            },
            // valueLoadString(v ref, b []byte)
            "syscall/js.valueLoadString": (sp) => {
              sp >>>= 0;
              let str = loadValue(sp + 8);
              loadSlice(sp + 16).set(str);
            },
            // func valueInstanceOf(v ref, t ref) bool
            "syscall/js.valueInstanceOf": (sp) => {
              sp >>>= 0, this.mem.setUint8(sp + 24, loadValue(sp + 8) instanceof loadValue(sp + 16) ? 1 : 0);
            },
            // func copyBytesToGo(dst []byte, src ref) (int, bool)
            "syscall/js.copyBytesToGo": (sp) => {
              sp >>>= 0;
              let dst = loadSlice(sp + 8), src = loadValue(sp + 32);
              if (!(src instanceof Uint8Array || src instanceof Uint8ClampedArray)) {
                this.mem.setUint8(sp + 48, 0);
                return;
              }
              let toCopy = src.subarray(0, dst.length);
              dst.set(toCopy), setInt64(sp + 40, toCopy.length), this.mem.setUint8(sp + 48, 1);
            },
            // func copyBytesToJS(dst ref, src []byte) (int, bool)
            "syscall/js.copyBytesToJS": (sp) => {
              sp >>>= 0;
              let dst = loadValue(sp + 8), src = loadSlice(sp + 16);
              if (!(dst instanceof Uint8Array || dst instanceof Uint8ClampedArray)) {
                this.mem.setUint8(sp + 48, 0);
                return;
              }
              let toCopy = src.subarray(0, dst.length);
              dst.set(toCopy), setInt64(sp + 40, toCopy.length), this.mem.setUint8(sp + 48, 1);
            },
            debug: (value) => {
              console.log(value);
            }
          }
        };
      }
      async run(instance) {
        if (!(instance instanceof WebAssembly.Instance))
          throw new Error("Go.run: WebAssembly.Instance expected");
        this._inst = instance, this.mem = new DataView(this._inst.exports.mem.buffer), this._values = [
          // JS values that Go currently has references to, indexed by reference id
          NaN,
          0,
          null,
          !0,
          !1,
          globalThis,
          this
        ], this._goRefCounts = new Array(this._values.length).fill(1 / 0), this._ids = /* @__PURE__ */ new Map([
          // mapping from JS values to reference ids
          [0, 1],
          [null, 2],
          [!0, 3],
          [!1, 4],
          [globalThis, 5],
          [this, 6]
        ]), this._idPool = [], this.exited = !1;
        let offset = 4096, strPtr = (str) => {
          let ptr = offset, bytes = encoder.encode(str + "\0");
          return new Uint8Array(this.mem.buffer, offset, bytes.length).set(bytes), offset += bytes.length, offset % 8 !== 0 && (offset += 8 - offset % 8), ptr;
        }, argc = this.argv.length, argvPtrs = [];
        this.argv.forEach((arg) => {
          argvPtrs.push(strPtr(arg));
        }), argvPtrs.push(0), Object.keys(this.env).sort().forEach((key) => {
          argvPtrs.push(strPtr(`${key}=${this.env[key]}`));
        }), argvPtrs.push(0);
        let argv = offset;
        if (argvPtrs.forEach((ptr) => {
          this.mem.setUint32(offset, ptr, !0), this.mem.setUint32(offset + 4, 0, !0), offset += 8;
        }), offset >= 12288)
          throw new Error("total length of command line and environment variables exceeds limit");
        this._inst.exports.run(argc, argv), this.exited && this._resolveExitPromise(), await this._exitPromise;
      }
      _resume() {
        if (this.exited)
          throw new Error("Go program has already exited");
        this._inst.exports.resume(), this.exited && this._resolveExitPromise();
      }
      _makeFuncWrapper(id) {
        let go = this;
        return function() {
          let event = { id, this: this, args: arguments };
          return go._pendingEvent = event, go._resume(), event.result;
        };
      }
    };
  })();

  // pkg/admin/src/rogueId.ts
  var RogueIdElement = class extends HTMLElement {
    tooltip;
    constructor() {
      super(), this.tooltip = document.createElement("span"), this.setupTooltip();
    }
    setupTooltip() {
      this.tooltip.style.position = "absolute", this.tooltip.style.display = "none", this.tooltip.style.backgroundColor = "black", this.tooltip.style.color = "white", this.tooltip.style.padding = "4px 8px", this.tooltip.style.borderRadius = "4px", this.tooltip.style.fontSize = "12px", this.tooltip.style.zIndex = "1000";
    }
    connectedCallback() {
      document.body.appendChild(this.tooltip), this.addEventListener("mouseover", this.handleMouseOver), this.addEventListener("mouseout", this.handleMouseOut), this.addEventListener("click", this.handleClick);
    }
    handleMouseOver = (event) => {
      let rogueId = this.getAttribute("data-rogue-id"), isDel = this.getAttribute("data-is-del");
      rogueId && (document.querySelectorAll(
        `[data-rogue-id="${rogueId}"]`
      ).forEach((el) => el.classList.add("userhovered")), isDel === "true" && (this.tooltip.style.backgroundColor = "#FFEBEB", this.tooltip.style.color = "#B30000"), this.tooltip.textContent = rogueId, this.tooltip.style.top = `${event.clientY - 40}px`, this.tooltip.style.display = "block", this.tooltip.style.left = `${event.clientX - 10}px`);
    };
    handleMouseOut = () => {
      let rogueId = this.getAttribute("data-rogue-id");
      rogueId && (document.querySelectorAll(
        `[data-rogue-id="${rogueId}"]`
      ).forEach((el) => el.classList.remove("userhovered")), this.tooltip.style.display = "none");
    };
    handleClick = (event) => {
      let rogueId = this.getAttribute("data-rogue-id");
      rogueId && (navigator.clipboard.writeText(rogueId), this.style.borderColor = "green", this.tooltip.innerHTML = "ID copied to clipboard");
    };
    disconnectedCallback() {
      this.tooltip && this.tooltip.parentElement && this.tooltip.parentElement.removeChild(this.tooltip);
    }
  };
  customElements.define("rogue-id", RogueIdElement);
})();
//# sourceMappingURL=main.js.map
