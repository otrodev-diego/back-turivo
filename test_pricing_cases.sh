#!/bin/bash

# Script para probar los 6 casos de validación del endpoint de pricing
# Basado en el documento Validacion_v2.2.0_Turivo_Pricing.md

# Configuración
API_BASE_URL="http://localhost:8083"
TOKEN="your-token-here"  # Reemplazar con token real

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Ejecutando Tests de Validación v2.2.0 - Turivo Pricing${NC}"
echo "=================================================="

# Función para hacer request y validar resultado
test_case() {
    local case_name="$1"
    local expected_final="$2"
    local expected_commission="$3"
    local expected_driver="$4"
    local json_data="$5"
    
    echo -e "\n${YELLOW}📋 Caso: $case_name${NC}"
    echo "Request: $json_data"
    
    # Hacer request
    response=$(curl -s -X POST "$API_BASE_URL/api/v1/pricing/quote" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "$json_data")
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Error en request${NC}"
        return 1
    fi
    
    echo "Response: $response"
    
    # Extraer valores (usando jq si está disponible)
    if command -v jq &> /dev/null; then
        final_fare=$(echo "$response" | jq -r '.finalFare // "null"')
        commission=$(echo "$response" | jq -r '.commission // "null"')
        driver=$(echo "$response" | jq -r '.driverPayout // "null"')
        
        echo "Final Fare: $final_fare"
        echo "Commission: $commission"
        echo "Driver Payout: $driver"
        
        # Validar tolerancia (±0.5%)
        if [ "$final_fare" != "null" ] && [ "$expected_final" != "null" ] && [ "$final_fare" != "" ] && [ "$expected_final" != "" ]; then
            # Usar awk para cálculo más robusto
            diff=$(awk "BEGIN {printf \"%.2f\", ($final_fare - $expected_final) / $expected_final * 100}")
            abs_diff=$(awk "BEGIN {printf \"%.2f\", ($diff < 0) ? -$diff : $diff}")
            
            if (( $(echo "$abs_diff <= 0.5" | bc -l 2>/dev/null || echo "0") )); then
                echo -e "${GREEN}✅ Aprobado - Δ% = $diff%${NC}"
            else
                echo -e "${RED}❌ Rechazado - Δ% = $diff% (tolerancia: ±0.5%)${NC}"
            fi
        else
            echo -e "${RED}❌ Error: No se pudo extraer valores de la respuesta${NC}"
        fi
    else
        echo -e "${YELLOW}⚠️  jq no disponible, no se puede validar automáticamente${NC}"
    fi
}

# Caso 1: T004 Transfer Aeropuerto
test_case "T004 Transfer Aeropuerto" "49140" "9828" "39312" '{
  "serviceCode": "T004",
  "distanceKm": 25,
  "vehicleTypeId": "van_premium",
  "segmentId": "B2B",
  "zoneId": "urbana",
  "scheduleId": "punta",
  "currencyCode": "CLP"
}'

# Caso 2: T003 Transfer Urbano
test_case "T003 Transfer Urbano" "13200" "2640" "10560" '{
  "serviceCode": "T003",
  "distanceKm": 10,
  "vehicleTypeId": "van_estandar",
  "segmentId": "B2C",
  "zoneId": "mixta",
  "scheduleId": "normal",
  "currencyCode": "CLP"
}'

# Caso 3: T014 Tour Viña del Mar
test_case "T014 Tour Viña del Mar" "420000" "84000" "336000" '{
  "serviceCode": "T014",
  "zoneId": "rural",
  "scheduleId": "normal",
  "currencyCode": "CLP"
}'

# Caso 4: T015 Tour Cajón del Maipo
test_case "T015 Tour Cajón del Maipo" "390000" "78000" "312000" '{
  "serviceCode": "T015",
  "zoneId": "interregional",
  "scheduleId": "punta",
  "currencyCode": "CLP"
}'

# Caso 5: Ruta Integrada 2 paradas + 1h espera
test_case "Ruta Integrada 2+1" "47920" "9584" "38336" '{
  "serviceCode": "T003",
  "distanceKm": 18,
  "vehicleTypeId": "sedan_ejecutivo",
  "segmentId": "B2C",
  "zoneId": "urbana",
  "scheduleId": "normal",
  "paradas": 2,
  "horasEspera": 1,
  "currencyCode": "CLP"
}'

# Caso 6: Ruta Integrada Mínimo activado
test_case "Ruta Integrada Mínimo" "36000" "7200" "28800" '{
  "serviceCode": "T009",
  "distanceKm": 5,
  "vehicleTypeId": "van_estandar",
  "segmentId": "B2B",
  "zoneId": "mixta",
  "scheduleId": "nocturno",
  "paradas": 0,
  "horasEspera": 0,
  "currencyCode": "CLP"
}'

echo -e "\n${BLUE}📊 Resumen de Tests Completados${NC}"
echo "=================================================="
echo -e "${GREEN}✅ Todos los casos han sido probados${NC}"
echo -e "${YELLOW}⚠️  Revisar los resultados arriba para validación${NC}"
