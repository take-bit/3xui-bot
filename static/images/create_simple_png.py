#!/usr/bin/env python3
"""
Простой скрипт для создания PNG файла без PIL
Создает минимальный PNG файл 400x300 с темно-синим фоном
"""

import struct
import zlib

def create_simple_png():
    # PNG signature
    png_signature = b'\x89PNG\r\n\x1a\n'
    
    # IHDR chunk (Image Header)
    width = 400
    height = 300
    bit_depth = 8
    color_type = 2  # RGB
    compression = 0
    filter_method = 0
    interlace = 0
    
    ihdr_data = struct.pack('>IIBBBBB', width, height, bit_depth, color_type, 
                           compression, filter_method, interlace)
    ihdr_crc = zlib.crc32(b'IHDR' + ihdr_data) & 0xffffffff
    ihdr_chunk = struct.pack('>I', 13) + b'IHDR' + ihdr_data + struct.pack('>I', ihdr_crc)
    
    # IDAT chunk (Image Data) - темно-синий фон с градиентом
    image_data = bytearray()
    
    for y in range(height):
        # Фильтр типа 0 (None)
        image_data.append(0)
        
        for x in range(width):
            # Создаем градиент от темно-синего к синему
            gradient_factor = (x + y) / (width + height)
            
            # Базовый темно-синий: (30, 58, 138)
            r = int(30 + gradient_factor * 50)  # 30-80
            g = int(58 + gradient_factor * 100)  # 58-158  
            b = int(138 + gradient_factor * 50)  # 138-188
            
            # Добавляем RGB пиксель
            image_data.extend([r, g, b])
    
    compressed_data = zlib.compress(image_data)
    
    idat_crc = zlib.crc32(b'IDAT' + compressed_data) & 0xffffffff
    idat_chunk = struct.pack('>I', len(compressed_data)) + b'IDAT' + compressed_data + struct.pack('>I', idat_crc)
    
    # IEND chunk
    iend_crc = zlib.crc32(b'IEND') & 0xffffffff
    iend_chunk = struct.pack('>I', 0) + b'IEND' + struct.pack('>I', iend_crc)
    
    # Собираем PNG
    png_data = png_signature + ihdr_chunk + idat_chunk + iend_chunk
    
    # Записываем файл
    with open('static/images/bot_banner.png', 'wb') as f:
        f.write(png_data)
    
    print("✅ Создан PNG баннер: static/images/bot_banner.png")
    print(f"   Размер: {width}x{height} пикселей")
    print(f"   Размер файла: {len(png_data)} байт")
    print("   Цвет: Градиент от темно-синего к синему")

if __name__ == "__main__":
    create_simple_png()
