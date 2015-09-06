for i in 1..100 do
    puts "#{i},#{i * 3},#{i * 1000}"
    STDOUT.flush
    sleep 0.3
end
