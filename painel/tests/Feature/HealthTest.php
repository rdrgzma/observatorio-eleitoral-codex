<?php

namespace Tests\Feature;

use Illuminate\Support\Facades\DB;
use Tests\TestCase;

class HealthTest extends TestCase
{
    public function test_root_declares_absence_of_models(): void
    {
        $this->get('/')->assertOk()->assertJsonPath('electoral_models', false);
    }

    public function test_framework_health(): void
    {
        $this->get('/up')->assertOk();
    }

    public function test_readiness_checks_database(): void
    {
        DB::shouldReceive('select')->once()->with('SELECT 1')->andReturn([(object) ['value' => 1]]);
        $this->get('/readyz')->assertOk()->assertJsonPath('database', 'reachable');
    }

    public function test_database_failure_is_unavailable_without_details(): void
    {
        DB::shouldReceive('select')->once()->andThrow(new \RuntimeException('sensitive connection details'));
        $this->get('/readyz')->assertStatus(503)->assertExactJson(['status' => 'unavailable']);
    }

    public function test_missing_route(): void
    {
        $this->getJson('/missing')->assertNotFound();
    }
}
