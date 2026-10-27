<?php

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Route;

Route::get('/', fn () => response()->json([
    'service' => 'painel',
    'status' => 'scaffold',
    'electoral_models' => false,
    'notice' => 'Plataforma em desenvolvimento; nenhum modelo eleitoral implementado.',
]));

// Local readiness check. Never returns connection details or credentials.
Route::get('/readyz', function () {
    try {
        DB::select('SELECT 1');
    } catch (Throwable $exception) {
        return response()->json(['status' => 'unavailable'], 503);
    }

    return response()->json(['status' => 'ok', 'database' => 'reachable']);
});
