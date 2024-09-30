import sqlite3

def save_to_database(data, db_path='textbook_data.db'):
    conn = sqlite3.connect(db_path)
    cursor = conn.cursor()
    
    cursor.execute('''
    CREATE TABLE IF NOT EXISTS content_items
    (id INTEGER PRIMARY KEY, type TEXT, content TEXT, tags TEXT)
    ''')
    
    for item in data:
        cursor.execute('''
        INSERT INTO content_items (type, content, tags)
        VALUES (?, ?, ?)
        ''', (item['type'], item['content'], ','.join(item['tags'])))
    
    conn.commit()
    conn.close()