/// <reference path="./wasm_exec.d.ts" />

// deno-lint-ignore no-sloppy-imports
import "./wasm_exec.js";

import mainwasm from "./mainwasm.ts";

import "jsr:@supabase/functions-js/edge-runtime.d.ts";

import { decodeBase64 } from "https://deno.land/std@0.224.0/encoding/base64.ts";

const go = new Go();

const module = decodeBase64(mainwasm);

const { instance } = await WebAssembly.instantiate(module, go.importObject);

go.run(instance);

// deno-lint-ignore no-explicit-any
const handle = (globalThis as any).handle as (req: Request) => Response;

Deno.serve(async (req) => {
  const body = await req.json();
  const res = handle(body);
  return res;
});
